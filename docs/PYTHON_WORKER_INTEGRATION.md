# Python Microservice Integration Guide

This guide shows how your Python microservice worker should interact with the translation job system.

## Database Models (SQLAlchemy example)

```python
from sqlalchemy import Column, String, Integer, DateTime, Text, ForeignKey
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship
from datetime import datetime
import uuid

class TranslationJob(Base):
    __tablename__ = 'translation_jobs'
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    novel_id = Column(UUID(as_uuid=True), nullable=False, index=True)
    from_lang = Column(String(10), nullable=False)
    target_lang = Column(String(10), nullable=False)
    status = Column(String(20), nullable=False, default='PENDING')
    progress = Column(Integer, nullable=False, default=0)
    total_subtasks = Column(Integer, nullable=False, default=0)
    completed_subtasks = Column(Integer, nullable=False, default=0)
    staging_prefix = Column(Text, nullable=True)
    error_message = Column(Text, nullable=True)
    created_by = Column(UUID(as_uuid=True), nullable=True)
    started_at = Column(DateTime, nullable=True)
    finished_at = Column(DateTime, nullable=True)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=False, default=datetime.utcnow, onupdate=datetime.utcnow)
    
    subtasks = relationship("TranslationSubtask", back_populates="job")

class TranslationSubtask(Base):
    __tablename__ = 'translation_subtasks'
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    job_id = Column(UUID(as_uuid=True), ForeignKey('translation_jobs.id', ondelete='CASCADE'), nullable=False)
    entity_type = Column(String(20), nullable=False)
    entity_id = Column(UUID(as_uuid=True), nullable=False)
    parent_volume_id = Column(UUID(as_uuid=True), nullable=True)
    seq = Column(Integer)
    priority = Column(Integer, nullable=False, default=100)
    status = Column(String(20), nullable=False, default='PENDING')
    result_path = Column(Text, nullable=True)
    result_text = Column(Text, nullable=True)
    error_message = Column(Text, nullable=True)
    started_at = Column(DateTime, nullable=True)
    finished_at = Column(DateTime, nullable=True)
    created_at = Column(DateTime, nullable=False, default=datetime.utcnow)
    updated_at = Column(DateTime, nullable=False, default=datetime.utcnow, onupdate=datetime.utcnow)
    
    job = relationship("TranslationJob", back_populates="subtasks")
```

## Worker Implementation

```python
import redis
import json
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from datetime import datetime

# Constants
ENTITY_TYPE_CHAPTER = "chapter"
ENTITY_TYPE_VOLUME = "volume"
ENTITY_TYPE_NOVEL = "novel"

STATUS_PENDING = "PENDING"
STATUS_IN_PROGRESS = "IN_PROGRESS"
STATUS_DONE = "DONE"
STATUS_FAILED = "FAILED"

class TranslationWorker:
    def __init__(self, db_url, redis_url):
        self.engine = create_engine(db_url)
        self.Session = sessionmaker(bind=self.engine)
        self.redis = redis.from_url(redis_url)
        
    def run(self):
        """Main worker loop"""
        print("Translation worker started...")
        
        while True:
            # Block and wait for a job from Redis queue
            _, job_data = self.redis.blpop("translation_queue")
            job_id = json.loads(job_data)["job_id"]
            
            print(f"Processing job: {job_id}")
            self.process_job(job_id)
    
    def process_job(self, job_id):
        """Process a translation job"""
        session = self.Session()
        
        try:
            # Update job status to IN_PROGRESS
            job = session.query(TranslationJob).filter_by(id=job_id).first()
            if not job:
                print(f"Job {job_id} not found")
                return
            
            job.status = STATUS_IN_PROGRESS
            job.started_at = datetime.utcnow()
            session.commit()
            
            # Process all subtasks
            while True:
                # Claim next pending subtask (with lock)
                subtask = session.query(TranslationSubtask)\
                    .filter_by(job_id=job_id, status=STATUS_PENDING)\
                    .order_by(TranslationSubtask.priority, TranslationSubtask.seq)\
                    .with_for_update(skip_locked=True)\
                    .first()
                
                if not subtask:
                    break  # All subtasks processed
                
                # Mark subtask in progress
                subtask.status = STATUS_IN_PROGRESS
                subtask.started_at = datetime.utcnow()
                session.commit()
                
                try:
                    # Process the subtask
                    self.process_subtask(session, subtask)
                    
                    # Mark subtask done
                    subtask.status = STATUS_DONE
                    subtask.finished_at = datetime.utcnow()
                    
                    # Update job progress
                    job.completed_subtasks += 1
                    job.progress = int((job.completed_subtasks / job.total_subtasks) * 100)
                    
                    session.commit()
                    print(f"Subtask {subtask.id} completed ({job.progress}%)")
                    
                except Exception as e:
                    # Mark subtask failed
                    subtask.status = STATUS_FAILED
                    subtask.error_message = str(e)
                    subtask.finished_at = datetime.utcnow()
                    
                    # Mark job failed
                    job.status = STATUS_FAILED
                    job.error_message = f"Subtask {subtask.id} failed: {str(e)}"
                    
                    session.commit()
                    print(f"Subtask {subtask.id} failed: {e}")
                    return
            
            # All subtasks done - promote to production
            print(f"All subtasks done for job {job_id}. Promoting...")
            self.promote_translations(session, job)
            
            # Mark job completed
            job.status = "COMPLETED"
            job.progress = 100
            job.finished_at = datetime.utcnow()
            session.commit()
            
            print(f"Job {job_id} completed successfully")
            
        except Exception as e:
            print(f"Job {job_id} failed: {e}")
            job.status = STATUS_FAILED
            job.error_message = str(e)
            session.commit()
        finally:
            session.close()
    
    def process_subtask(self, session, subtask):
        """Process a single subtask"""
        entity_type = subtask.entity_type
        entity_id = subtask.entity_id
        
        if entity_type == ENTITY_TYPE_CHAPTER:
            self.translate_chapter(session, subtask)
        elif entity_type == ENTITY_TYPE_VOLUME:
            self.translate_volume(session, subtask)
        elif entity_type == ENTITY_TYPE_NOVEL:
            self.translate_novel(session, subtask)
        else:
            raise ValueError(f"Unknown entity type: {entity_type}")
    
    def translate_chapter(self, session, subtask):
        """Translate a chapter"""
        # 1. Fetch chapter and translation
        chapter = session.query(Chapter).filter_by(id=subtask.entity_id).first()
        chapter_trans = session.query(ChapterTranslation)\
            .filter_by(chapter_id=chapter.id)\
            .filter_by(lang=subtask.job.from_lang)\
            .first()
        
        if not chapter_trans:
            raise ValueError(f"No translation found for chapter {chapter.id}")
        
        # 2. Translate title and content
        translated_title = self.translate_text(chapter_trans.title, subtask.job.target_lang)
        translated_content = self.translate_text(chapter_trans.content, subtask.job.target_lang)
        
        # 3. Store in staging (S3 or file system)
        staging_data = {
            "title": translated_title,
            "content": translated_content
        }
        staging_path = self.save_to_staging(subtask.job.id, subtask.id, staging_data)
        
        # 4. Update subtask with result path
        subtask.result_path = staging_path
    
    def translate_volume(self, session, subtask):
        """Translate volume metadata"""
        # 1. Fetch volume translation
        volume_trans = session.query(VolumeTranslation)\
            .filter_by(volume_id=subtask.entity_id)\
            .filter_by(lang=subtask.job.from_lang)\
            .first()
        
        if not volume_trans:
            raise ValueError(f"No translation found for volume {subtask.entity_id}")
        
        # 2. Translate title and description
        translated_title = self.translate_text(volume_trans.title, subtask.job.target_lang)
        translated_desc = self.translate_text(volume_trans.description, subtask.job.target_lang) if volume_trans.description else None
        
        # 3. Store small text directly in result_text
        subtask.result_text = json.dumps({
            "title": translated_title,
            "description": translated_desc
        })
    
    def translate_novel(self, session, subtask):
        """Translate novel metadata"""
        # Similar to volume translation
        novel_trans = session.query(NovelTranslation)\
            .filter_by(novel_id=subtask.entity_id)\
            .filter_by(lang=subtask.job.from_lang)\
            .first()
        
        if not novel_trans:
            raise ValueError(f"No translation found for novel {subtask.entity_id}")
        
        translated_title = self.translate_text(novel_trans.title, subtask.job.target_lang)
        translated_desc = self.translate_text(novel_trans.description, subtask.job.target_lang) if novel_trans.description else None
        
        subtask.result_text = json.dumps({
            "title": translated_title,
            "description": translated_desc
        })
    
    def translate_text(self, text, target_lang):
        """Call your ML translation model"""
        # TODO: Implement actual translation
        # This is where you call your translation model
        # For now, just return placeholder
        return f"[{target_lang}] {text}"
    
    def save_to_staging(self, job_id, subtask_id, data):
        """Save data to staging storage (S3 or local)"""
        # TODO: Implement S3 upload or file system storage
        # Return the path/URL to the staged file
        path = f"staging/{job_id}/{subtask_id}.json"
        # Save to S3 or local file...
        return path
    
    def promote_translations(self, session, job):
        """Atomically promote all staging results to production tables"""
        print("Starting promotion transaction...")
        
        try:
            # Start explicit transaction
            session.begin_nested()
            
            # Get all subtasks for this job
            subtasks = session.query(TranslationSubtask)\
                .filter_by(job_id=job.id)\
                .all()
            
            for subtask in subtasks:
                if subtask.entity_type == ENTITY_TYPE_CHAPTER:
                    # Load from staging
                    staging_data = self.load_from_staging(subtask.result_path)
                    
                    # Upsert chapter translation
                    existing = session.query(ChapterTranslation)\
                        .filter_by(chapter_id=subtask.entity_id, lang=job.target_lang)\
                        .first()
                    
                    if existing:
                        existing.title = staging_data["title"]
                        existing.content = staging_data["content"]
                        existing.updated_at = datetime.utcnow()
                    else:
                        new_trans = ChapterTranslation(
                            chapter_id=subtask.entity_id,
                            lang=job.target_lang,
                            title=staging_data["title"],
                            content=staging_data["content"]
                        )
                        session.add(new_trans)
                
                elif subtask.entity_type == ENTITY_TYPE_VOLUME:
                    result_data = json.loads(subtask.result_text)
                    
                    existing = session.query(VolumeTranslation)\
                        .filter_by(volume_id=subtask.entity_id, lang=job.target_lang)\
                        .first()
                    
                    if existing:
                        existing.title = result_data["title"]
                        existing.description = result_data.get("description")
                        existing.updated_at = datetime.utcnow()
                    else:
                        new_trans = VolumeTranslation(
                            volume_id=subtask.entity_id,
                            lang=job.target_lang,
                            title=result_data["title"],
                            description=result_data.get("description")
                        )
                        session.add(new_trans)
                
                elif subtask.entity_type == ENTITY_TYPE_NOVEL:
                    result_data = json.loads(subtask.result_text)
                    
                    existing = session.query(NovelTranslation)\
                        .filter_by(novel_id=subtask.entity_id, lang=job.target_lang)\
                        .first()
                    
                    if existing:
                        existing.title = result_data["title"]
                        existing.description = result_data.get("description")
                        existing.updated_at = datetime.utcnow()
                    else:
                        new_trans = NovelTranslation(
                            novel_id=subtask.entity_id,
                            lang=job.target_lang,
                            title=result_data["title"],
                            description=result_data.get("description")
                        )
                        session.add(new_trans)
            
            # Commit all changes atomically
            session.commit()
            print("Promotion completed successfully")
            
        except Exception as e:
            session.rollback()
            raise Exception(f"Promotion failed: {e}")
    
    def load_from_staging(self, path):
        """Load data from staging"""
        # TODO: Load from S3 or file system
        return {"title": "...", "content": "..."}


if __name__ == "__main__":
    worker = TranslationWorker(
        db_url="postgresql://user:pass@localhost/dbname",
        redis_url="redis://localhost:6379/0"
    )
    worker.run()
```

## How to Push Job to Redis from Go

In your Go service after creating the job:

```go
import (
    "context"
    "encoding/json"
    "github.com/redis/go-redis/v9"
)

func (s *TranslationJobService) pushToQueue(ctx context.Context, jobID string) error {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    payload, _ := json.Marshal(map[string]string{"job_id": jobID})
    return rdb.RPush(ctx, "translation_queue", payload).Err()
}
```

Add this call after creating the job in `CreateTranslationJob` method.

## Environment Setup

Python requirements:
```
sqlalchemy
psycopg2-binary
redis
boto3  # if using S3
```

Docker Compose additions:
```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
  
  translation-worker:
    build: ./python-worker
    environment:
      - DATABASE_URL=postgresql://user:pass@postgres:5432/dbname
      - REDIS_URL=redis://redis:6379/0
      - S3_BUCKET=translation-staging  # optional
    depends_on:
      - postgres
      - redis
```

## Testing

1. Create a job via API:
```bash
curl -X POST http://localhost:8080/api/translation-jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"novel_id": "uuid", "target_lang": "en"}'
```

2. Check job status:
```bash
curl http://localhost:8080/api/translation-jobs/{job_id} \
  -H "Authorization: Bearer $TOKEN"
```

3. Monitor worker logs to see progress
