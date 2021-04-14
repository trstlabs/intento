<template>
	<div class="sp-component sp-transfer-list" >
        <p>sgsfgdss</p>
		<div class="sp-transfer-list__header sp-component-title">
			<h3>Transactions</h3>
			<span>|</span>
			<span>A list of your recent transactions</span>
		</div>
		<table
			class="sp-transfer-list__table sp-box sp-shadow"
			v-if="bankAddress && transactions.length > 0"
		>
			<thead>
				<tr>
					<th class="sp-transfer-list__status">STATUS</th>
					<th class="sp-transfer-list__table__address">ADDRESS / DETAILS</th>
					<th class="sp-transfer-list__table__amount">AMOUNT</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="(tx, index) in sentTransactions" v-bind:key="index">
					<p> TEST </p>
					<td class="sp-transfer-list__status">
						<div class="sp-transfer-list__status__wrapper">
							<div
								
							>
								<span
									
								/>
							</div>
							<div class="sp-transfer-list__status__action">
								<div class="sp-transfer-list__status__action__text">
									{{ getTxText(tx) }}
								</div>
								<!--<div class="sp-transfer-list__status__action__date">
									{{ getFmtTime(tx.response.timestamp) }}
								</div>-->
							</div>
						</div>
					</td>
					<td class="sp-transfer-list__table__address"><p> TTTEST </p>
						{{ getTxDetails(tx) }}
					</td>
					<td
						class="sp-transfer-list__table__amount"
						v-if="
							tx.body.messages[0]['@type'] == '/cosmos.bank.v1beta1.MsgSend'
						"
					>
						<div
							v-for="(token, index) in tx.body.messages[0].amount"
							v-bind:key="'am' + index"
						>
							{{
								tx.body.messages[0].from_address == bankAddress
									? '-' + token.amount + ' ' + token.denom.toUpperCase()
									: '+' + token.amount + ' ' + token.denom.toUpperCase()
							}}
						</div>
					</td>
					<td
						class="sp-transfer-list__table__amount"
						v-else-if="
							tx.body.messages[0]['@type'] ==
							'/ibc.applications.transfer.v1.MsgTransfer'
						"
					>
						<div>
							{{
								tx.body.messages[0].sender == bankAddress
									? '-' +
									  tx.body.messages[0].token.amount +
									  ' ' +
									  tx.body.messages[0].token.denom.toUpperCase()
									: '+' +
									  tx.body.messages[0].token.amount +
									  ' ' +
									  tx.body.messages[0].token.denom.toUpperCase()
							}}
						</div>
					</td>
					<td
						class="sp-transfer-list__table__amount"
						v-else-if="
							tx.body.messages[0]['@type'] ==
							'/ibc.core.channel.v1.MsgRecvPacket'
						"
					>
						<div>
							{{
								getDecoded(tx.body.messages[0].packet.data).receiver ==
								bankAddress
									? '+' +
									  getDecoded(tx.body.messages[0].packet.data).amount +
									  ' IBC/' +
									  tx.body.messages[0].packet.destination_port.toUpperCase() +
									  '/' +
									  tx.body.messages[0].packet.destination_channel.toUpperCase() +
									  '/' +
									  getDecoded(
											tx.body.messages[0].packet.data
									  ).denom.toUpperCase()
									: '-' +
									  getDecoded(tx.body.messages[0].packet.data).amount +
									  ' IBC/' +
									  tx.body.messages[0].packet.destination_port.toUpperCase() +
									  '/' +
									  tx.body.messages[0].packet.destination_channel.toUpperCase() +
									  '/' +
									  getDecoded(
											tx.body.messages[0].packet.data
									  ).denom.toUpperCase()
							}}
						</div>
					</td>
					<td class="sp-transfer-list__table__amount" v-else></td>
				</tr>
			</tbody>
		</table>

		<table class="sp-transfer-list__table sp-box sp-shadow" v-else>
			<tbody>
				<tr>
					<td class="sp-transfer-list__status">
						<div class="sp-transfer-list__status__wrapper">
							<div
								class="sp-transfer-list__status__icon sp-transfer-list__status__icon__empty"
							>
								<span class="sp-icon sp-icon-Transactions" />
							</div>
							<div class="sp-transfer-list__status__action">
								<div class="sp-transfer-list__status__action__text">
									No transactions yet
								</div>
								<div
									class="sp-transfer-list__status__action__date"
									v-if="!bankAddress"
								>
									Add or unlock a wallet to see recent transactions
								</div>
							</div>
						</div>
					</td>
					<td class="sp-transfer-list__table__address"></td>
					<td class="sp-transfer-list__table__amount"></td>
				</tr>
			</tbody>
		</table>
	</div>
</template>
<script>
import axios from 'axios'

export default {
	name: 'SpTransferList',
//	props: { address: String, refresh: Boolean },
	data: function () {
		return {
			bankAddress: '',
            GetTxsEvent: {},
        _Subscriptions: new Set(),
		sentTransactions: {},
		receivedTransactions: {}
		}
	},
	computed: {
	
	
		fullBalances() {
			return this.balances.map((x) => {
				this.addMapping(x)
				return x
			})
		},
		transactions() {
		
		let sent =
				this.sentTransactions.txs?.map((tx, index) => {
					tx = this.sentTransactions[index]
		
					return tx
				}) ?? []
							
				console.log(sent)
				let received =
				this.receivedTransactions.txs?.map((tx, index) => {
					tx.response = this.receivedTransactions.tx_result[index]
					return tx
				}) ?? []
			return [...sent, ...received].sort(
				(a, b) => b.response.height - a.response.height
			)
		
				
		}
	},
	beforeCreate() {
		
	},

   

	async created() {
	
			this.bankAddress = this.$store.state.account.address
			if (this.bankAddress != '') {
            try {
				let sent = (await axios.get(process.env.VUE_APP_RPC + '/tx_search?query=' + '"message.sender%3D%27' + this.bankAddress + '%27"')).data;                
            
				this.sentTransactions = JSON.stringify(sent.result)
				    console.log(this.sentTransactions) 
					 console.log(sent.result) 

              let received = (await axios.get(process.env.VUE_APP_RPC + '/tx_search?query=' + '"message.recipient%3D%27' + this.bankAddress + '%27"')).data;                
              
				this.receivedTransactions = JSON.stringify(received.result)
				  console.log(this.receivedTransactions) }

			 catch (e) {
                //console.error(new SpVuexError('QueryClient:ServiceGetTxsEvent', 'API Node Unavailable. Could not perform query.'));
           console.log("ERROR" + e) 
	
			 } }},
	methods: {

       
         async ServiceGetTxsEvent({ commit }, { ...key }) {
            try {
                let params=Object.values(key)
                let value = (await axios.get(process.env.VUE_APP_RPC + '/tx_search?query=' + '"message.sender%3D%27' + this.bankAddress + '%27"')).data;                
                console.log(value) 
                /*
                while (all && value.pagination && value.pagination.next_key!=null) {
                    let next_values=(await (await initQueryClient(rootGetters)).queryPostAll.apply(null,[...params, {'pagination.key':value.pagination.next_key}] )).data;
                    
                    for (let prop of Object.keys(next_values)) {
                        if (Array.isArray(next_values[prop])) {
                            value[prop]=[...value[prop], ...next_values[prop]]
                        }else{
                            value[prop]=next_values[prop]
                        }
                    }
                    console.log(value)
                }
                */
                this.QUERY({query: 'GetTxsEvent', key, value });
                if (subscribe)
                    commit('SUBSCRIBE', { action: 'ServiceGetTxsEvent', payload: key });
            }
            catch (e) {
                //console.error(new SpVuexError('QueryClient:ServiceGetTxsEvent', 'API Node Unavailable. Could not perform query.'));
           console.log("ERROR" + e) }
        },

async QUERY( {query, key, value }) {
            this[query][JSON.stringify(key)] = value;
        },
        
		getFmtTime(time) {
			const momentTime = dayjs(time)
			return momentTime.format('D MMM, YYYY')
		},
		getDecoded(packet) {
			try {
				return JSON.parse(decode(packet))
			} catch (e) {
				return {}
			}
		},
		getTxText(tx) {
			let text = ''
			if (tx.code != 0) {
				text = '(Failed) '
			}
			if (tx.body.messages.length > 1) {
				text = text + 'Multiple messages'
			} else {
				if (
					tx.body.messages[0]['@type'] == '/cosmos.bank.v1beta1.MsgSend' ||
					tx.body.messages[0]['@type'] ==
						'/ibc.applications.transfer.v1.MsgTransfer'
				) {
					if (tx.body.messages[0].from_address == this.bankAddress) {
						text = text + 'Sent to'
					}
					if (tx.body.messages[0].to_address == this.bankAddress) {
						text = text + 'Received from'
					}
					if (tx.body.messages[0].sender == this.bankAddress) {
						text = text + 'IBC Sent to'
					}
					if (tx.body.messages[0].receiver == this.bankAddress) {
						text = text + 'IBC Received from'
					}
				} else {
					let packet
					switch (tx.body.messages[0]['@type']) {
						case '/ibc.core.channel.v1.MsgChannelOpenAck':
							text = text + 'IBC Channel Open Ack'
							break
						case '/ibc.core.channel.v1.MsgChannelOpenConfirm':
							text = text + 'IBC Channel Open Confirm'
							break
						case '/ibc.core.channel.v1.MsgChannelOpenTry':
							text = text + 'IBC Channel Open Try'
							break
						case '/ibc.core.channel.v1.MsgRecvPacket':
							packet = this.getDecoded(tx.body.messages[0].packet.data)
							if (packet.sender == this.bankAddress) {
								text = text + 'IBC Sent to'
							} else {
								if (packet.receiver == this.bankAddress) {
									text = text + 'IBC Received from'
								} else {
									text = text + 'IBC Recv Packet'
								}
							}
							break
						case '/ibc.core.channel.v1.MsgAcknowledgement':
							text = text + 'IBC Ack Packet'
							break
						case '/ibc.core.channel.v1.MsgTimeout':
							text = text + 'IBC Timeout Packet'
							break
						case '/ibc.core.channel.v1.MsgChannelOpenInit':
							text = text + 'IBC Channel Open Init'
							break
						case '/ibc.core.client.v1.MsgCreateClient':
							text = text + 'IBC Client Create'
							break
						case '/ibc.core.client.v1.MsgUpdateClient':
							text = text + 'IBC Client Update'
							break
						case '/ibc.core.connection.v1.MsgConnectionOpenAck':
							text = text + 'IBC Connection Open Ack'
							break
						case '/ibc.core.connection.v1.MsgConnectionOpenInit':
							text = text + 'IBC Connection Open Init'
							break
						case '/ibc.core.connection.v1.MsgConnectionOpenConfirm':
							text = text + 'IBC Connection Open Confirm'
							break
						case '/ibc.core.connection.v1.MsgConnectionOpenTry':
							text = text + 'IBC Connection Open Try'
							break
						default:
							text = text + 'Message'
							break
					}
				}
			}
			return text
		},
		getTxDetails(tx) {
			let text = ''
			if (tx.body.messages.length > 1) {
				text = text + '-'
			} else {
				if (
					tx.body.messages[0]['@type'] == '/cosmos.bank.v1beta1.MsgSend' ||
					tx.body.messages[0]['@type'] ==
						'/ibc.applications.transfer.v1.MsgTransfer'
				) {
					if (tx.body.messages[0].from_address == this.bankAddress) {
						text = text + tx.body.messages[0].to_address
					}
					if (tx.body.messages[0].to_address == this.bankAddress) {
						text = text + tx.body.messages[0].from_address
					}
					if (tx.body.messages[0].sender == this.bankAddress) {
						let chain = this.$store.getters['common/relayers/chainFromChannel'](
							tx.body.messages[0].source_channel
						)
						text = text + chain + ':' + tx.body.messages[0].receiver
					}
					if (tx.body.messages[0].receiver == this.bankAddress) {
						let chain = this.$store.getters['common/relayers/chainToChannel'](
							tx.body.messages[0].source_channel
						)
						text = text + chain + ':' + tx.body.messages[0].receiver
					}
				} else {
					let packet
					switch (tx.body.messages[0]['@type']) {
						case '/ibc.core.channel.v1.MsgChannelOpenAck':
							text =
								text +
								tx.body.messages[0].port_id +
								' / ' +
								tx.body.messages[0].channel_id
							break
						case '/ibc.core.channel.v1.MsgChannelOpenConfirm':
							text =
								text +
								tx.body.messages[0].port_id +
								' / ' +
								tx.body.messages[0].channel_id
							break
						case '/ibc.core.channel.v1.MsgChannelOpenTry':
							text =
								text +
								tx.body.messages[0].port_id +
								' / ' +
								tx.body.messages[0].previous_channel_id +
								' / ' +
								tx.body.messages[0].counterparty_version
							break
						case '/ibc.core.channel.v1.MsgRecvPacket':
							packet = this.getDecoded(tx.body.messages[0].packet.data)
							if (packet.sender == this.bankAddress) {
								text = text + 'IBC:' + packet.receiver
							} else {
								if (packet.receiver == this.bankAddress) {
									text = text + 'IBC:' + packet.sender
								} else {
									text = text + 'IBC Recv Packet'
								}
							}
							break
						case '/ibc.core.channel.v1.MsgAcknowledgement':
							text =
								text +
								tx.body.messages[0].packet.source_port +
								':' +
								tx.body.messages[0].packet.source_channel +
								' <-> ' +
								tx.body.messages[0].packet.destination_port +
								':' +
								tx.body.messages[0].packet.destination_channel
							break
						case '/ibc.core.channel.v1.MsgTimeout':
							text = text + 'IBC Timeout Packet'
							break
						case '/ibc.core.channel.v1.MsgChannelOpenInit':
							text = text + tx.body.messages[0].port_id
							break
						case '/ibc.core.client.v1.MsgCreateClient':
							text = text + tx.body.messages[0].signer
							break
						case '/ibc.core.client.v1.MsgUpdateClient':
							text = text + tx.body.messages[0].client_id
							break
						case '/ibc.core.connection.v1.MsgConnectionOpenAck':
							text =
								text +
								tx.body.messages[0].connection_id +
								' / ' +
								tx.body.messages[0].counterparty_connection_id
							break
						case '/ibc.core.connection.v1.MsgConnectionOpenInit':
							text = text + tx.body.messages[0].client_id
							break
						case '/ibc.core.connection.v1.MsgConnectionOpenConfirm':
							text = text + tx.body.messages[0].connection_id
							break
						case '/ibc.core.connection.v1.MsgConnectionOpenTry':
							text =
								text +
								tx.body.messages[0].client_id +
								' / ' +
								tx.body.messages[0].previous_connection_id
							break
						default:
							text = text + 'Message'
							break
					}
				}
			}
			return text
		}
	},

}
</script>
