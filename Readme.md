# Nextcloud Uploader Action (Go 1.23)

This action uploads a file to a Nextcloud server.

## Usage

```yaml
- name: Nextcloud Upload
  uses: eskopp/nextcloud-upload-file@main
  with:
    file-path: './path/to/file'
    nextcloud-url: 'https://nextcloud.example.com/remote.php/webdav/'
    username: ${{ secrets.NEXTCLOUD_USERNAME }}
    password: ${{ secrets.NEXTCLOUD_PASSWORD }}
```
