Start dev service

```bash
cd server
hzn dev service start -S -v             
```

Place you `.wav` files under `/sets/<set-name>` folder and start the test:

```bash
./test-client -b localhost:1883 -u fft-client -p client-pass -c fft-test --result_topic results
```
