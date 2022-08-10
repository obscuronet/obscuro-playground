const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idGenerateViewingKey = "generateViewingKey";
const idStatus = "status";
const pathGenerateViewingKey = "/generateviewingkey/";
const pathSubmitViewingKey = "/submitviewingkey/";
const methodPost = "post";
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};
const metamaskRequestAccounts = "eth_requestAccounts";
const metamaskPersonalSign = "personal_sign";
const personalSignPrefix = "vk";

const initialize = () => {
    const generateViewingKeyButton = document.getElementById(idGenerateViewingKey);
    const statusArea = document.getElementById(idStatus);

    generateViewingKeyButton.addEventListener(eventClick, async () => {
        if (typeof ethereum === "undefined") {
            statusArea.innerText = "`ethereum` object is not available. Please install and enable MetaMask."
            return
        }

        const accounts = await ethereum.request({method: metamaskRequestAccounts});
        if (accounts.length === 0) {
            statusArea.innerText = "No MetaMask accounts found."
            return
        }
        // Accounts is "An array of a single, hexadecimal Ethereum address string.", so we grab the single entry at index zero.
        const account = accounts[0];

        const addressJson = {"address": account}
        const viewingKeyResp = await fetch(
            pathGenerateViewingKey, {
                method: methodPost,
                headers: jsonHeaders,
                body: JSON.stringify(addressJson)
            }
        );
        if (!viewingKeyResp.ok) {
            statusArea.innerText = "Failed to generate viewing key."
            return
        }

        const viewingKey = await viewingKeyResp.text();

        const signature = await ethereum.request({
            method: metamaskPersonalSign,
            // Without a prefix such as 'vk', personal_sign transforms the data for security reasons.
            params: [personalSignPrefix + viewingKey, account]
        }).catch(_ => { return -1 })
        if (signature === -1) {
            statusArea.innerText = "Failed to sign viewing key."
            return
        }

        const signedViewingKeyJson = {"signature": signature}
        const submitViewingKeyResp = await fetch(
            pathSubmitViewingKey, {
                method: methodPost,
                headers: jsonHeaders,
                body: JSON.stringify(signedViewingKeyJson)
            }
        );
        if (submitViewingKeyResp.ok) {
            statusArea.innerText = `Account: ${account}\nViewing key: ${viewingKey}\nSigned bytes: ${signature}`
        } else {
            statusArea.innerText = "Failed to submit viewing key to enclave."
        }
    })
}

window.addEventListener(eventDomLoaded, initialize);