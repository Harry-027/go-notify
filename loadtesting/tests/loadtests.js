import http from "k6/http";
import { check, sleep } from "k6";
import { Counter } from "k6/metrics";

let ErrorCount = new Counter("errors");

export const options = {
    vus: 10,
    duration: "15s",
    thresholds: {
        errors: ["count<10"]
    }
};

export default function() {
    var hostip = 'your machine ip' // where the api server is listening
    var port = ':3001'
    var hostAddr = `${hostip}${port}`
    var commonParam = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    for (var id = 1; id <= 100; id++) {
        var rnd = Math.random().toString(36).replace(/[^a-z]+/g, '')
        var signupUrl = `http://${hostAddr}/auth/signup`;
        var mail = `aaa${id}${rnd}@gmail.com`
        var signupPayload = JSON.stringify({
            email: mail,
            password: 'bbb',
            confirm_password: 'bbb',
            role: 'admin'
        });
        let signupRes = http.post(signupUrl, signupPayload, commonParam);
        let signupSuccess = check(signupRes, {
            "status is 200": r => {
                console.log("Signup Status code ::", r.status)
                return r.status === 200
            }
        });

        if (!signupSuccess) {
            ErrorCount.add(1);
        }

        // login process ::
        var token = ''
        var loginurl = `http://${hostAddr}/auth/login`;
        var loginpayload = JSON.stringify({
            email: mail,
            password: 'bbb',
        });
        let loginRes = http.post(loginurl, loginpayload, commonParam);
        token = loginRes.body.token
        let loginSuccess = check(loginRes, {
            "status is 200": r => {
                console.log("Login Status code ::", r.status)
                return r.status === 200
            }
        });

        if (!loginSuccess) {
            ErrorCount.add(1);
        }

        // Fetch the users ::
        var fetchUserUrl = `http://${hostAddr}/api/users`;
        var userParam = {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': token
            },
        };
        let userRes = http.get(fetchUserUrl, userParam);
        let userSuccess = check(loginRes, {
            "status is 200": r => {
                console.log("User Status code ::", r.status)
                return r.status === 200
            }
        });

        if (!userSuccess) {
            ErrorCount.add(1);
        }


    }
    sleep(2);
}