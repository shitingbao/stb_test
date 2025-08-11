
import json
import requests

payload={
    "code": 200,
    "taskId": 17,
    "fileId": 23,
    "labels": [
        {
            "imageId": 69,
            "label": "light_zz_wz",
            "type": "rect",
            "x": 0.40859375,
            "y": 0.5173611111111112,
            "w": 0.11875,
            "h": 0.020833333333333332
        },
        {
            "imageId": 69,
            "label": "zz_zx_wz",
            "type": "rect",
            "x": 0.502734375,
            "y": 0.5034722222222222,
            "w": 0.07109375,
            "h": 0.16527777777777777
        },
        {
            "imageId": 70,
            "label": "dzsd_m",
            "type": "rect",
            "x": 0.940234375,
            "y": 0.6958333333333333,
            "w": 0.11796875,
            "h": 0.36666666666666664
        }
    ],
    "msg": ""
}

headers = {
            "Content-Type": "application/json"
        }
print("payload", payload)
server_url = "http://127.0.0.1:8808/api/v1/start_autolabel/back"
response = requests.request("POST", server_url, headers = headers, data = json.dumps(payload))