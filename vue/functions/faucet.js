//const restEndpoint = process.env.NODE_URL
//const restEndpoint = "https://node.trustpriceprotocol.com"
const ADDRESS_PREFIX = 'cosmos';
import { SigningStargateClient } from "@cosmjs/stargate";
import { DirectSecp256k1HdWallet } from '@cosmjs/proto-signing';

const {Bech32} = require("@cosmjs/encoding")

const GoogleRecaptcha = require('google-recaptcha')
console.log(process.env.GOOGLE)
const googleRecaptcha = new GoogleRecaptcha({
  secret: process.env.GOOGLE
})





let signer

exports.handler = async function (event, context) {

  
  //signer = await Secp256k1Wallet.fromMnemonic(process.env.MNEMONIC)
  // console.log({signer})
  // let headers = {
  //   'Access-Control-Allow-Origin': '*',
  //   'Access-Control-Allow-Methods': 'GET,POST,PUT,DELETE,OPTIONS',
  //   'Access-Control-Allow-Headers':
  //     'Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With'
  // }
  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(
    process.env.MNEMONIC,
    makeCosmoshubPath(0),
    ADDRESS_PREFIX
  )

  if (event.httpMethod === 'POST') {
    if (event.body) {
      let body = JSON.parse(event.body)
      console.log({body})
      let recipient = body.recipient
      try {
        Bech32.decode(recipient)
      } catch (err) {
        console.log("NEED TO RESOLVE")
        console.log({err})

        try {
          recipient = await getValue(recipient)
          Bech32.decode(recipient)
        } catch (err) {
          console.log("STILL NEED TO RESOLVE")
          return {
            statusCode: 400,
            body: JSON.stringify(err)
          }
        }
      }


      let recaptchaResponse = body.recaptchaToken
      let response
      try  {
        response = await new Promise((resolve, reject) => {
          console.log({googleRecaptcha})
          googleRecaptcha.verify({ response: recaptchaResponse }, async (error, response) => {
            if (error) { reject(error) } else { resolve(response) }
          })
        })
        console.log({response})
      } catch (error) {
        console.log({error})
         return {
          
          statusCode: 400,
          body: error.message
        }
      }
      
      if (!response.success) {
        console.error(response)
        return {
          statusCode: 400,
          body: JSON.stringify(response)
        }
      } else {
        try {
          const result = await submitWithCosmJS(recipient)
          return {
            statusCode: 200,
            body: JSON.stringify(result.data)
          }
        } catch (error) {
          console.log({error})
          return handleAxiosError(error)
        }
      }
    } else {
      return {
        statusCode: 404,
        body: '¯\\_(ツ)_/¯'
      }
    }
  } else {
    return {
      statusCode: 200,
      body: ':)'
    }
  }
}

async function submitWithCosmJS(recipient) {
  const [firstAccount] = await wallet.getAccounts();

const rpcEndpoint = 'https://cli.trustpriceprotocol.com';
const client = await SigningStargateClient.connectWithSigner(rpcEndpoint, wallet);

const amount = {
  denom: "tpp",
  amount: "5",
};

const result = await client.sendTokens(firstAccount.address, recipient, [amount], "Welcome to the Trust Price Protocol community");
assertIsBroadcastTxSuccess(result);

  //const [{address}] = await signer.getAccounts();
  //console.log({address})
 /* const memo = ''
  const msg = {
    type:  'cosmos-sdk/MsgSend',
    value: {
        amount:  [{ amount: '5', denom: 'tpp' }],
        from_address: address,
        to_address: recipient
    }
  };
  const fee = {
    amount: coins(0, ''),
    gas: "200000",
  };
  console.log({restEndpoint})
  const client = new SigningCosmosClient(restEndpoint, address, signer);
  console.log({client})
  return client.signAndPost([msg], fee, memo);*/

}

function handleAxiosError (error) {
  console.error(error)
  return {
    statusCode: !error.response ? 500 : error.response.status,
    body: !error.response ? error.message : error.response.statusText + (error.response.data && error.response.data.error ? '\n' + error.response.data.error : '')
  }
}