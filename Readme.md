# Nextcloud Uploader Action (Go 1.23)

This action uploads a file to a Nextcloud server.

## Usage

```yaml
- name: Nextcloud Upload
  uses: eskopp/nextcloud-upload-file@v0.0.1
  with:
    file-path: "./path/to/file.txt"
    nextcloud-url: "https://nextcloud.example.com/remote.php/webdav/"
    username: ${{ secrets.NEXTCLOUD_USERNAME }}
    password: ${{ secrets.NEXTCLOUD_PASSWORD }}
    override: "false" # optional (If the file exists and override is not set or is false, the server returns error code 204.)
    rename: "false" # optional (Replace false with the name to which the file is to be renamed. If rename is false or not there, the file will not be renamed.)
```
