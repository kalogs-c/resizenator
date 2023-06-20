# Resizenator

Resizenator is a powerful cloud function designed specifically for Google Cloud
Functions and Google Cloud Storage. It efficiently monitors the specified Google
Cloud Storage bucket for new file uploads and automatically resizes the images
according to the configurations provided in the `config.yaml` file. With
Resizenator, you can easily resize images into multiple sizes with customizable
options.

## Features

- Monitors Google Cloud Storage for new file uploads.
- Resizes uploaded images into multiple sizes.
- Customizable options to configure the resizing process.

## Configuration

Resizenator offers a flexible configuration system to tailor the resizing
process to your specific needs. Adjust the parameters in the `config.yaml` file
based on your requirements. The following configuration options are available:

- **max_resizing_concurrency**: Sets the maximum number of concurrent image
  resizing operations. Adjust this value based on available resources and
  performance requirements.

- **prefix**: Specifies the prefix to be added to the resized images. Use this
  to differentiate between the original and resized versions.

- **delete_after_upload**: Controls whether the original image should be deleted
  after resizing and uploading. Set this to `true` to remove the original image
  once the resizing process is complete.

- **sizes**: Defines the dimensions (in pixels) to which the images should be
  resized. Resizenator generates resized versions for each size provided.

- **interpolation_algorithm**: Specifies the interpolation algorithm to be used
  during the resizing process. Choose from available options: BiLinear,
  NearestNeighbor, CatmullRom, or ApproxBiLinear. Select the algorithm based on
  the desired image quality and performance.

- **target_format**: Specifies the target image format for the resized images.
  Choose from available options: webp, jpeg/jpg, png, or gif. Resizenator
  converts the images to the specified format during the resizing process. If no
  format is specified, the original image format is used.

Before deploying Resizenator, ensure that you modify the `config.yaml` file to
align the resizing process with your requirements.

## Deployment

To deploy Resizenator, follow these steps:

1. Clone the Resizenator repository.
2. Configure the `config.yaml` file with your desired options.
3. Deploy the cloud function to Google Cloud Functions using the provided
   deployment instructions.
4. Set up appropriate permissions for the cloud function to access your Google
   Cloud Storage bucket.
5. Monitor the specified Google Cloud Storage bucket for new file uploads, and
   Resizenator will automatically resize the images according to your
   configurations.

## Contributions

Contributions to Resizenator are welcome! If you encounter any issues or have
ideas for improvements, please submit them through the issue tracker of the
Resizenator repository on GitHub. Additionally, you can fork the repository,
make changes, and create a pull request to contribute directly.
