function main(params) {
  return new Promise((resolve, reject) => {
    if (!params.messages || !params.messages[0] || !params.messages[0].value) {
      reject("Invalid arguments. Must include 'messages' JSON array with 'value' field");
    }
    const msgs = params.messages;
    const cats = [];
    for (let i = 0; i < msgs.length; i++) {
      const msg = msgs[i];
      for (let j = 0; j < msg.value.cats.length; j++) {
        const cat = msg.value.cats[j];
        console.log(`A ${cat.color} cat named ${cat.name} was received.`);
        cats.push(cat);
      }
    }
    resolve({
      cats,
    });
  });
}

exports.main = main;
