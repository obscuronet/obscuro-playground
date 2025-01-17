import { ethers } from "hardhat";
import { ManagementContract } from "../../typechain-types";

const setChallengePeriod = async function (mgmtContractAddr: string, challengPeriod: number) {
    const managementContract = await ethers.getContractAt(
        "ManagementContract",
        mgmtContractAddr
    ) as ManagementContract;


    console.log(`Setting challenge period to: ${challengPeriod}`);
    const tx = await managementContract.SetChallengePeriod(BigInt(challengPeriod));
    await tx.wait();
    console.log(`Successfully set challenge period to: ${challengPeriod}`);
    
    const mgmtContractChallengePeriod = await managementContract.GetChallengePeriod();
    if (BigInt(challengPeriod) !== mgmtContractChallengePeriod) {
        throw new Error(`Failed to set the challenge period to: ${challengPeriod}. Returned value is: ${mgmtContractChallengePeriod}`);
    }
}

const mgmtContractAddr = process.env.MGMT_CONTRACT_ADDRESS;
const challengePeriod = process.env.L1_CHALLENGE_PERIOD ?
    Number(process.env.L1_CHALLENGE_PERIOD) : 0;


if (!mgmtContractAddr) {
    console.error("Missing required environment variables: MGMT_CONTRACT_ADDRESS");
    process.exit(1);
}

setChallengePeriod(mgmtContractAddr, challengePeriod)
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

export default setChallengePeriod;