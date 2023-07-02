# Lilly Project Plan #

1. Add pkg/transform/edge_detect.go
   1. Display Edge Detection Dialog
      1. I would like a down-sampled display of edge detection with current values in the dialog.
         1. Selector for number of parallel channels
         2. Any other tuneable parameters?
      2. Apply edge detection convolution to loaded image, load resulting image.
2. Add pkg/transform/scale.go
   1. Display Scale Dialog
   2. Scale loaded image, load resulting scaled image.
3. Add pkg/transform/translate.go
4. Add pkg/transform/rotate.go