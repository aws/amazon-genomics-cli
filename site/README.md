## Amazon Genomics CLI Website
The Amazon Genomics CLI website is built with [hugo](https://gohugo.io/) and the [docsy](https://www.docsy.dev/) theme.

### Pre-requisites

1. Make sure that you have installed the [hugo](https://gohugo.io/getting-started/installing/) binary. If you're using MacOS you can use brew:

   ```bash
   $ brew install hugo
   ```

2. Pull down the docsy theme. The docsy theme submodule is listed under the `.gitmodules` file.

   ```bash
   $ cd site/
   $ git submodule update --init --recursive
   ```

   You may additionally have to pull down the CSS and font libraries required by the Docsy theme:

   ```bash
   $ git submodule update --init --recursive --depth 1
   ```

4. Install the [CSS processing libraries](https://www.docsy.dev/docs/getting-started/#install-postcss) required by Docsy:

   ```bash
   $ npm install
   ```

### Releasing the latest docs to GitHub pages
Once we are ready to release a new version of the docs, you can run `make build-docs`.
Alternatively, you can:

```bash
$ cd site/
$ hugo
$ cd ..
```

This will update the documentation under the `docs/` directory. Afterwards, Create a new PR with the changes:

```bash
$ git add docs/
$ git commit -m "docs: update website for agc vX.Y.Z"
$ git push <remote> docs
```

### Developing locally

From the root of the repository, run `make start-docs`.  
Alternatively, you can:

```bash
$ cd site/
$ hugo server -D
```

Then you should be able to access the website at [http://localhost:1313/](http://localhost:1313/). With the `-D` flag
set, Hugo will render "draft" documents. To exclude these then remove the `-D` flag.

#### Adding new content
Follow the Docsy [guidance](https://www.docsy.dev/docs/adding-content/content/) for adding new content.

Content should be added in `content/<lang>/docs` under the appropriate language folder (e.g `content/en/docs` for English language docs.)

#### Styling content

Docsy integrates by default with [bootstrap4](https://getbootstrap.com/docs/4.0/getting-started/introduction/), so we
can leverage any of the classes available there.

If you'd like to override any class that docsy itself generates, add the scss file under `assets/scss/`.
You can find which files are available under `themes/docsy/assets/scss`.

### Command Reference Generation

Cobra can automatically generate markdown reference for all commands. If `agc` is run in the root of this project and
the folder `site/content/en/docs/Reference` is present then markdown will be generated (or updated) automatically. 