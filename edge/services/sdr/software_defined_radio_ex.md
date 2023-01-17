# Software-defined radio edge processing

This example uses software-defined radio (SDR) as an example of edge processing. With SDR, you can send raw data across the full radio spectrum to a cloud server for processing. The edge node processes the data locally and then sends less volume of more valuable data to a cloud processing service for extra processing.

This diagram shows the architecture for this SDR example:

![Example architecture](../../evtstreams/sdr2evtstreams/08_sdrarch.svg)

SDR edge processing is a fully featured example that consumes radio station audio, extracts speech, and converts the extracted speech into text. The example completes sentiment analysis on the text and makes the data and results available through a user interface where you can view the details of the data for each edge node. Use this example to learn more about edge processing.

SDR receives radio signals by using the digital circuitry in a computer CPU to handle the work to require a set of specialized analog circuitry. That analog circuitry is usually restricted by the breadth of radio spectrum it can receive. An analog radio receiver built to receive FM radio stations, for example, cannot receive radio signals from anywhere else on the radio spectrum. SDR can access large portions of the spectrum. If you do not have the SDR hardware, that you can use mock data. When you are using mock data, the audio from the internet stream is treated as though it was broadcast over FM and received on your edge node.

Before performing this task, register and unregister your edge device by performing the steps in [Installing the agent](https://open-horizon.github.io/docs/installing/registration.html).

This code contains these primary components.

|Component|Description|
|---------|-----------|
|[sdr service](https://github.com/open-horizon/examples/tree/master/edge/services/sdr)|Lower-level service accesses the hardware on the edge node|
|[ssdr2evtstreams service](https://github.com/open-horizon/examples/tree/master/edge/evtstreams/sdr2evtstreams)|Higher-level service receives data from the lower-level sdr service and completes local analysis of the data on the edge node. The sdr2evtstreams service then sends the processed data to the cloud back-end software.|

## Registering your device

Though this service can be run by using mock data on any edge device, if you are using an edge device like a Raspberry Pi with the SDR hardware, configure a kernel module to support your SDR hardware first. You must manually configure this module. Docker containers can establish a different distribution of Linux in their contexts, but the container cannot install kernel modules.

Complete these steps to configure this module:

1. As a root user, create a file that is named `/etc/modprobe.d/rtlsdr.conf`

   ```bash
   sudo nano /etc/modprobe.d/rtlsdr.conf
   ```

2. Add the following lines to the file:

   ```bash
   blacklist rtl2830
   blacklist rtl2832
   blacklist dvb_usb_rtl28xxu
   ```

3. Save the file and then restart before you continue:

   ```bash
   sudo reboot
   ```

4. Set the following Apache Kafka API key in your environment. This key is created for use with this example and is used for feeding the processed data that is gathered by your edge node to the IBM software-defined radio UI.

   ```bash
   export EVTSTREAMS_API_KEY=X2e8cSjbDAMk-ztJLaoi3uffy8qsQTnZttUjcHCfm7cp
   export EVTSTREAMS_BROKER_URL=broker-3-y420pyyyvhhmttz0.kafka.svc01.us-south.eventstreams.cloud.ibm.com:9093,broker-5-y420pyyyvhhmttz0.kafka.svc01.us-south.eventstreams.cloud.ibm.com:9093,broker-4-y420pyyyvhhmttz0.kafka.svc01.us-south.eventstreams.cloud.ibm.com:9093,broker-1-y420pyyyvhhmttz0.kafka.svc01.us-south.eventstreams.cloud.ibm.com:9093,broker-0-y420pyyyvhhmttz0.kafka.svc01.us-south.eventstreams.cloud.ibm.com:9093,broker-2-y420pyyyvhhmttz0.kafka.svc01.us-south.eventstreams.cloud.ibm.com:9093
   ```

5. To run the sdr2evtstreams service example on your edge node, you must register your edge node with the IBM/pattern-ibm.sdr2evtstreams deployment pattern. Perform the steps in [Preconditions for Using the SDR To IBM Event Streams Example Edge Service](https://github.com/open-horizon/examples/tree/master/edge/evtstreams/sdr2evtstreams).

6. Check the example web UI to see whether your edge node is sending results.

## SDR implementation details

### sdr low-level service

The lowest level of the software stack for this service includes the `sdr` service implementation. This service accesses local software-defined radio hardware by using the popular `librtlsdr` library and the derived `rtl_fm` and `rtl_power` utilities along with the `rtl_rpcd` daemon. For more information about the `librtlsdr` library, see [librtlsdr](https://github.com/librtlsdr/librtlsdr).

The `sdr` service directly controls the software-defined radio hardware to tune the hardware to a particular frequency to receive transmitted data, or to measure the signal strength across a specified spectrum. A typical workflow for the service can be to tune to a particular frequency to receive data from the station at that frequency. Then, the service can process the collected data.

### sdr2evtstreams high-level service

The `sdr2evtstreams` high-level service implementation uses both the `sdr` service REST API and the `gps` service REST API over the local private virtual Docker network. The `sdr2evtstreams` service receives data from the `sdr` service, and completes some local inference on the data to select the best stations for speech. Then, the `sdr2evtstreams` service uses Kafka to publish audio clips to the cloud by using Apache Kafka.

## What to do next

If you want to deploy your own software to an edge node, you must create your own edge services, and associated deployment pattern or deployment policy. For more information, see [Developing an edge service for devices](https://open-horizon.github.io/docs/developing/developing.html).
