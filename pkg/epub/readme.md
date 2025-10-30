
## 🧩 The return type: `*epub.EpubContent`

When your `EpubService.UploadAndParseEpub` finishes, it returns a **fully parsed representation** of the EPUB file in this structure:

```go
type EpubContent struct {
	Metadata     OPFMetadata
	Manifest     []OPFManifestItem
	Spine        []OPFItemRef
	ContentFiles map[string]ContentFile
	RawFiles     map[string][]byte
	OPFPath      string
}
```

That means your `EpubService` gives you a **complete, in-memory representation of the EPUB book** — including metadata, chapters, manifest entries, and raw file contents.

Now let’s unpack each part.

---

## 🧱 1️⃣ `Metadata` (book info)

```go
Metadata OPFMetadata
```

This comes from the `content.opf` file — the EPUB’s metadata section.

It contains things like:

| Field           | Example                          | Meaning           |
| --------------- | -------------------------------- | ----------------- |
| `Title`       | `["The Great Novel"]`          | Book title        |
| `Creator`     | `["John Doe"]`                 | Author            |
| `Language`    | `["en"]`                       | Original language |
| `Description` | `["A fantasy story about..."]` | Book summary      |
| `Publisher`   | `["LightNovel Press"]`         | Publisher name    |
| `Subject`     | `["Fantasy", "Adventure"]`     | Tags or genres    |
| `Date`        | `["2024-02-20"]`               | Publication date  |
| `Identifier`  | `["isbn:1234567890"]`          | Unique identifier |
| `Rights`      | `["All rights reserved"]`      | Copyright notice  |

👉 **You read this if you need the metadata of the book** — for example, when importing a novel into your database.

---

## 🧱 2️⃣ `Manifest` (the list of files inside EPUB)

```go
Manifest []OPFManifestItem
```

Each manifest item describes one file in the EPUB package, like:

```go
{
	ID: "chapter1",
	Href: "Text/chapter1.xhtml",
	MediaType: "application/xhtml+xml",
}
```

This is how the EPUB tells you **what files exist and what they contain** (HTML chapters, CSS, images, etc.).

👉 You use this if you need to:

* Extract the list of chapters
* Know which files are text, images, or stylesheets
* Resolve links inside the EPUB

---

## 🧱 3️⃣ `Spine` (the reading order)

```go
Spine []OPFItemRef
```

Each `ItemRef` points to a file ID in the `Manifest`.

This defines **the order of chapters** to display or process.

Example:

```go
Spine: [
  { IDRef: "chapter1" },
  { IDRef: "chapter2" },
  { IDRef: "chapter3" },
]
```

👉 You use this to know **which chapter comes next** and build your “reading flow”.

---

## 🧱 4️⃣ `ContentFiles` (parsed chapter contents)

```go
ContentFiles map[string]ContentFile
```

Each entry contains a chapter (or any readable file):

```go
type ContentFile struct {
	Path      string
	RawHTML   string
	PlainText string
	MediaType string
}
```

Example:

```go
ContentFiles["Text/chapter1.xhtml"] = ContentFile{
	Path: "Text/chapter1.xhtml",
	RawHTML: "<html><body><h1>Chapter 1</h1><p>Once upon a time...</p></body></html>",
	PlainText: "Chapter 1\nOnce upon a time...",
	MediaType: "application/xhtml+xml",
}
```

👉 This is the **main part you’ll want to read** —

it contains the **actual story text** (each chapter’s content).

You can use `PlainText` to get clean, readable text

or `RawHTML` if you want to preserve styling and markup.

---

## 🧱 5️⃣ `RawFiles` (all binary files)

```go
RawFiles map[string][]byte
```

Contains every file extracted from the EPUB — images, fonts, HTML, XML, etc.

Useful if you want to:

* Save cover images
* Inspect or export resources
* Debug EPUB structure

Example:

```go
RawFiles["Images/cover.jpg"] = []byte{0xFF, 0xD8, 0xFF, ...}
```

---

## 🧱 6️⃣ `OPFPath`

```go
OPFPath string
```

This tells you where the EPUB’s metadata file (`content.opf`) was found inside the ZIP — useful for resolving relative paths.

---

## 💡 In summary

| Field            | Description                               | When to use                               |
| ---------------- | ----------------------------------------- | ----------------------------------------- |
| `Metadata`     | Book info (title, author, language, etc.) | Display or import book metadata           |
| `Manifest`     | List of all files                         | To locate text/image files                |
| `Spine`        | Reading order                             | To load chapters in correct order         |
| `ContentFiles` | Parsed text content (HTML/plaintext)      | ✅**Read chapters or extract text** |
| `RawFiles`     | Every raw file (binary)                   | Save or debug EPUB content                |
| `OPFPath`      | Internal metadata path                    | Technical, rarely used directly           |

---

### ✅ **So, to answer your question directly:**

> “Which one is the content that I need to read?”

👉 You should read from:

```go
epubContent.ContentFiles
```

That’s where all the **actual chapters or story text** are.

For example:

```go
for path, content := range epubContent.ContentFiles {
	fmt.Println("File:", path)
	fmt.Println("Plain text:", content.PlainText)
}
```

That will print out each chapter’s content.

---

Would you like me to show you how to **extract only the readable chapters** (in correct order from `Spine`) into a single combined text or slice? That’s usually the next step in using EPUB data.
