# Ten Doc Site

This is the Ten Doc Site and it looks like [this](https://docs.ten.xyz/).

## Adding New Doc Site Pages

1. Clone this repository: https://github.com/ten-protocol/go-ten
2. Create your new content as a Markdown file in the `/docs` folder of the repo. Take care with the folder structure. 
   As a general rule, new titles in the left hand navigation menu should have their content contained in a separate 
   subfolder under docs, for example, `/docs/testnet` contains all the Markdown files relation to the testnet docs.
3. To have this new content shown in the left-hand navigation menu you need to modify the file 
   `/docs/_data/navigation.yml`. Follow the same format to add new headings and content titles. Remember to specify the 
   file type as `.html` for your new Markdown files, not `.md` when providing the URL.
4. Push your changes to tennet/go-ten
5. GitHub Pages will trigger a GitHub Action to use a Jekyll build job to create the static content and then publish 
   the pages at the custom URL.
6. Browse to https://docs.ten.xyz/ and check your content. Remember your browser will cache some of the pages so hit 
   refresh a few times if it looks like the content is missing or the navigation menu is incorrect.

## Updating Existing Doc Site Pages

This is the same as `Adding New Doc Site Pages`, above, but omitting step #3.

## Generating the HTML files locally

You can check your changes (e.g. the Markdown formatting) by serving the changes locally as follows:

1. [Install Jekyll](https://jekyllrb.com/docs/installation/)
2. Run the Jekyll `serve` command below to serve the updated pages on a local webserver at `http://127.0.0.1:4000`. Any 
   changes will be automatically reflected

   ```
   cd ./docs
   bundle exec jekyll serve
   ```
