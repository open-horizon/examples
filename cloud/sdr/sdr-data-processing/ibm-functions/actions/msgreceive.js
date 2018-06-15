function main(params) {
  var date = new Date();
  if (params.messages && params.messages[0] && params.messages[0].value) {
    console.log(params.messages.length + " messages received at: "+date.toLocaleString());
  } else {
    console.log("received message(s) at: "+date.toLocaleString());
  }
  return new Promise((resolve, reject) => {
    if (!params.messages || !params.messages[0] || !params.messages[0].value) {
      reject("Invalid arguments. Must include 'messages' JSON array with 'value' field");
    }
    const msgs = params.messages;
    //const cats = [];
    for (let i = 0; i < msgs.length; i++) {
      const msg = msgs[i];
      console.log("msg "+i+": "+msg.value);
      /* parse the msg json
      for (let j = 0; j < msg.value.cats.length; j++) {
        const cat = msg.value.cats[j];
        console.log(`A ${cat.color} cat named ${cat.name} was received.`);
        cats.push(cat);
      }
      */
    }
    resolve({
        "result": "Success: Message from IBM Message Hub processed."
        //cats,
    });
  });
}

exports.main = main;
