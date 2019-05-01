from flask import Flask, request
from flask_restful import Resource, Api
import json

app = Flask(__name__)
api = Api(app)

stopwordlist = ['the', 'a', 'for', 'of', 'on', 'at', 'as', 'are', 'am', 'is', 'before', 'but', 'do', 'by', 'in', 'with']

class CleanedText (Resource):
    def post(self):
        request_json = json.dumps(request.json)
        print("request json: %s"%(request_json))
        text = request.json['text']
        cleaned_text = ""
        if text:
            cleaned_text = ' '.join(w for w in text.split() if w not in stopwordlist)

        result = {"result": cleaned_text}
        print("result: %s"%(result))
        return result



api.add_resource(CleanedText, "/remove_stopword")


if __name__ == '__main__':   #app.run(host='0.0.0.0', port='5002')
	app.run(host='0.0.0.0', port='5002')

