# Horizon Monitor Service

Monitor is a dependent service of the [Horizon object cetection and classification example edge services](../visual_inferencing/README.md). Monitor is a tiny web server, implemented with Python Flask for monitoring the example's output. This service is not required by the visual inferencing services, but it enables a quick local check of these examples. When you are running these examples you can navigate to the host's port `5200` using your browser to see live output. There you should see output similar to this:

![example-page](https://raw.githubusercontent.com/MegaMosquito/achatina/master/art/page.png)

