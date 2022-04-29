# steganography-website
Steganography website application for hiding data inside images.

## API
The api can be accessed with the following routes:
- /api/encode
    - takes a picture and data value using the multipart form encoding.
  The picture must be a valid PNG, with the data being the text to hide in the
  picture file.
- /api/decode
    - takes a picture and data value using the multipart form encoding.
  The picture must be a valid PNG, with hidden text that will be sent back.
