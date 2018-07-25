# SDR Data Processing Component

This component processes all of the audio files that get sent to msg hub from the edge. The IBM Functions
actions in the `ibm-functions` dir use Watson services to convert the audio to text, and then do
sentiment analysis on it, and store the results (insights) in our compose postgres DB.

## Sample Watson Services Client

To experiment with the Watson services, edit `demo-processing/main.go` and `watson/stt/stt.go` and then run it using:
```
make stt-sample
```