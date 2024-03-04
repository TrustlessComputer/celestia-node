// Start a Express server at port 22260
const lighthouse = require("@lighthouse-web3/sdk");
const express = require('express');
const app = express();
app.use(express.json({limit: '50mb'}));
app.use(express.urlencoded({limit: '50mb'}));
const port = 22260;
const bodyParser = require('body-parser');
const axios = require('axios');

const NAMESPACE_FILECOIN = "tcfilecoin"

app.use(express.static('public'));
app.use(bodyParser.json());

app.post(`/filecoin/store`, async (req, res) => {
   // wrap in try catch block
    try {
        const data = req.body.data;
        if (!data) {
            return res.status(400).send('Invalid data');
        }
        // base64 decode to bytes arrays from req.body.data string
        const bytes = Buffer.from(data, 'base64');

        // write to tmp file with fileName is timestamp in milliseconds
        const fileName = new Date().getTime().toString();
        const filePath = `/tmp/${fileName}`;
        require('fs').writeFileSync(filePath, bytes);
        const dealParams = {
            network: 'calibration',
        };
        let resp = undefined;
        if (process.env.api_env === "mainnet") {
             resp = await lighthouse.upload( `/tmp/${fileName}`, getLighthouseAPIKey());
        } else {
            resp = await lighthouse.upload( `/tmp/${fileName}`, getLighthouseAPIKey(), false, dealParams);
        }
        console.log("resp from upload: ", resp);

        await sleep(5000);

        let params = {
            cid: resp.data.Hash,
            network: "testnet",
        };

        if (process.env.api_env === "mainnet") {
            params.network = "mainnet";
        }

        let response = await axios.get("https://api.lighthouse.storage/api/lighthouse/get_proof", {
            params: params,
        });
        const dealStatusResp = await  lighthouse.dealStatus(resp.data.Hash);
        let tcFileHash = resp.data.Hash;
        if (response.data.dealInfo && response.data.dealInfo.length > 0) {
            tcFileHash += '_' + response.data.dealInfo[0].dealId;
        }
        const result = {
            resp,
            cid: resp.data.Hash,
            dealStatus: dealStatusResp.data,
            proofResp: response.data,
            tcFileHash: tcFileHash,
        };
        return res.status(200).send(result);
    } catch (error) {
        return res.status(500).send(error);
    }
});

app.get(`/filecoin/get/tcfilecoin/:fileHash`, async (req, res) => {
    console.log("get file");
    try {
        const fileHash = req.params.fileHash;
        if (!fileHash) {
            return res.status(400).send('Invalid fileHash');
        }
        // split fileHash by _
        const fileHashArr = fileHash.split('_');
        if (fileHashArr.length < 1) {
            return res.status(400).send('Invalid fileHash');
        }

        const lighthouseDealDownloadEndpoint = 'https://gateway.lighthouse.storage/ipfs/'
        let response = await axios({
            method: 'GET',
            url: `${lighthouseDealDownloadEndpoint}${fileHashArr[0]}`,
            responseType: 'stream',
        });

        // how to response to the client (postman) ? Copilot please help me
        response.data.pipe(res);
        // return res.status(200).write(response);
    } catch (error) {
        return res.status(500).send(error);
    }
});

function sleep(ms) {
    return new Promise((resolve) => {
        setTimeout(resolve, ms);
    });
}

// getConfig return Lighthouse API key
function getLighthouseAPIKey() {
    let API_KEY = "11a8de93.35d708ea8ff547bca116f2f7519b052f";
    if (process.env.LIGHTHOUSE_API_KEY) {
        API_KEY = process.env.LIGHTHOUSE_API_KEY;
    }
    return API_KEY;
}


app.listen(port, () => {
    console.log(`Server is running at http://localhost:${port}`);
});
