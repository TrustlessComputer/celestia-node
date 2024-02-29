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
            num_copies: 2,
            repair_threshold: 28800,
            renew_threshold: 240,
            miner: ["t017840"],
            network: 'calibration',
            add_mock_data: 2
        };
        const resp = await lighthouse.upload( `/tmp/${fileName}`, getLighthouseAPIKey());

        console.log("resp", resp);
        await sleep(100000);


        let response = await axios.get("https://api.lighthouse.storage/api/lighthouse/get_proof", {
            params: {
                cid: resp.data.Hash,
                // network: "testnet" // Change the network to mainnet when ready
            }
        })
        console.log("response", response);

        const dealID = response.data.deal_id;

        console.log(resp);

        return res.status(200).send(cid);
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
