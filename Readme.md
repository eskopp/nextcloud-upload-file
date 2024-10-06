# Nextcloud Uploader Action (Go 1.23)

This action uploads a file to a Nextcloud server.

## Usage

```yaml
- name: Nextcloud Upload
  uses: eskopp/nextcloud-upload-file@v0.0.4
  with:
    file-path: "./path/to/file.txt"
    nextcloud-url: "https://nextcloud.example.com/remote.php/webdav/"
    username: ${{ secrets.NEXTCLOUD_USERNAME }}
    password: ${{ secrets.NEXTCLOUD_PASSWORD }}
    override: "false" # optional (If the file exists and override is not set or is false, the server returns error code 204.)
    rename: "false" # optional (Replace false with the name to which the file is to be renamed. If rename is empty or not there, the file will not be renamed.)
    zip: "false" # optional (Creates a zip archive with the name of the file)
    date: "true"  # optional (creates a date stamp in the format YYYY_MM_DD)
    time: "true" # optional (Creates a date stamp in the format HH_MM_SS)
```
