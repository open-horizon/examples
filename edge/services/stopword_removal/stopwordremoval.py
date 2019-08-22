from wsgiref.simple_server import make_server
import json

stopwordlist = ['the', 'a', 'for', 'of', 'on', 'at', 'as', 'are', 'am', 'is', 'before', 'but', 'do', 'by', 'in', 'with']

def clean_text(environ, start_response):
    path = environ.get('PATH_INFO')
    request_body_size = int(environ.get('CONTENT_LENGTH', 0))

    data = environ['wsgi.input'].read(request_body_size)

    string = json.loads(data)
    text = string['text']

    cleaned_text = ""
    if text:
        cleaned_text = ' '.join(w for w in text.split() if w not in stopwordlist)

    send = json.dumps({"result": cleaned_text})
    response_body = send.encode('utf-8')
    print(response_body)

    status = "200 OK"
    response_headers = [('Content-Type', 'application/json'), ("Content-Length", str(len(response_body)))]
    start_response(status, response_headers)
    return [response_body]

print("stropwordremoval running...")

httpd = make_server(
    '0.0.0.0', 5002, clean_text)

httpd.serve_forever()

