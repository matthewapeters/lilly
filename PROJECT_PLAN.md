# Lilly Project Plan #

1. Add pkg/transform/edge_detect.go
   1. Display Edge Detection Dialog (initial done)
      1. I would like a down-sampled display of edge detection with current values in the dialog. (initial done - not scaled)
         1. Selector for number of parallel channels (NA)
         2. Any other tuneable parameters? (F, S, Threshold - done)
      2. Apply edge detection convolution to loaded image, load resulting image. (initial done)
      3. Phase 2: add histogram of luminance to aid in threshold setting. (done)
      4. Remove Test button, allow histogram and edge-enchance to be triggered on changes to T, S.
      5. Allow edge-detection in layers and selected windows to enable special-treatment enhancement to parts of the image
2. Add pkg/transform/scale.go
   1. Display Scale Dialog
   2. Scale loaded image, load resulting scaled image.
3. Add pkg/transform/translate.go
4. Add pkg/transform/rotate.go
5.