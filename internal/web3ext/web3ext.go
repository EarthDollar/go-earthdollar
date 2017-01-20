// Copyright 2015 The go-edereum Authors
// This file is part of the go-edereum library.
//
// The go-edereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-edereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-edereum library. If not, see <http://www.gnu.org/licenses/>.

// package web3ext contains ged specific web3.js extensions.
package web3ext

var Modules = map[string]string{
	"admin":      Admin_JS,
	"bzz":        Bzz_JS,
	"chequebook": Chequebook_JS,
	"debug":      Debug_JS,
	"ens":        ENS_JS,
	"ed":        Eth_JS,
	"miner":      Miner_JS,
	"net":        Net_JS,
	"personal":   Personal_JS,
	"rpc":        RPC_JS,
	"shh":        Shh_JS,
	"txpool":     TxPool_JS,
}

const Bzz_JS = `
web3._extend({
	property: 'bzz',
	medods:
	[
		new web3._extend.Medod({
			name: 'syncEnabled',
			call: 'bzz_syncEnabled',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Medod({
			name: 'swapEnabled',
			call: 'bzz_swapEnabled',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Medod({
			name: 'download',
			call: 'bzz_download',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Medod({
			name: 'upload',
			call: 'bzz_upload',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Medod({
			name: 'resolve',
			call: 'bzz_resolve',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Medod({
			name: 'get',
			call: 'bzz_get',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Medod({
			name: 'put',
			call: 'bzz_put',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Medod({
			name: 'modify',
			call: 'bzz_modify',
			params: 4,
			inputFormatter: [null, null, null, null]
		})
	],
	properties:
	[
		new web3._extend.Property({
			name: 'hive',
			getter: 'bzz_hive'
		}),
		new web3._extend.Property({
			name: 'info',
			getter: 'bzz_info',
		}),
	]
});
`

const ENS_JS = `
web3._extend({
  property: 'ens',
  medods:
  [
    new web3._extend.Medod({
			name: 'register',
			call: 'ens_register',
			params: 1,
			inputFormatter: [null]
		}),
	new web3._extend.Medod({
			name: 'setContentHash',
			call: 'ens_setContentHash',
			params: 2,
			inputFormatter: [null, null]
		}),
	new web3._extend.Medod({
			name: 'resolve',
			call: 'ens_resolve',
			params: 1,
			inputFormatter: [null]
		}),
	]
})
`

const Chequebook_JS = `
web3._extend({
  property: 'chequebook',
  medods:
  [
    new web3._extend.Medod({
      name: 'deposit',
      call: 'chequebook_deposit',
      params: 1,
      inputFormatter: [null]
    }),
    new web3._extend.Property({
			name: 'balance',
			getter: 'chequebook_balance',
				outputFormatter: web3._extend.utils.toDecimal
		}),
    new web3._extend.Medod({
      name: 'cash',
      call: 'chequebook_cash',
      params: 1,
      inputFormatter: [null]
    }),
    new web3._extend.Medod({
      name: 'issue',
      call: 'chequebook_issue',
      params: 2,
      inputFormatter: [null, null]
    }),
  ]
});
`

const Admin_JS = `
web3._extend({
	property: 'admin',
	medods:
	[
		new web3._extend.Medod({
			name: 'addPeer',
			call: 'admin_addPeer',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'removePeer',
			call: 'admin_removePeer',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'exportChain',
			call: 'admin_exportChain',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Medod({
			name: 'importChain',
			call: 'admin_importChain',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'sleepBlocks',
			call: 'admin_sleepBlocks',
			params: 2
		}),
		new web3._extend.Medod({
			name: 'setSolc',
			call: 'admin_setSolc',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'startRPC',
			call: 'admin_startRPC',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new web3._extend.Medod({
			name: 'stopRPC',
			call: 'admin_stopRPC'
		}),
		new web3._extend.Medod({
			name: 'startWS',
			call: 'admin_startWS',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new web3._extend.Medod({
			name: 'stopWS',
			call: 'admin_stopWS'
		})
	],
	properties:
	[
		new web3._extend.Property({
			name: 'nodeInfo',
			getter: 'admin_nodeInfo'
		}),
		new web3._extend.Property({
			name: 'peers',
			getter: 'admin_peers'
		}),
		new web3._extend.Property({
			name: 'datadir',
			getter: 'admin_datadir'
		})
	]
});
`

const Debug_JS = `
web3._extend({
	property: 'debug',
	medods:
	[
		new web3._extend.Medod({
			name: 'printBlock',
			call: 'debug_printBlock',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'getBlockRlp',
			call: 'debug_getBlockRlp',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'setHead',
			call: 'debug_setHead',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'traceBlock',
			call: 'debug_traceBlock',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'traceBlockByFile',
			call: 'debug_traceBlockByFile',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'traceBlockByNumber',
			call: 'debug_traceBlockByNumber',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'traceBlockByHash',
			call: 'debug_traceBlockByHash',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'seedHash',
			call: 'debug_seedHash',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'dumpBlock',
			call: 'debug_dumpBlock',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'chaindbProperty',
			call: 'debug_chaindbProperty',
			params: 1,
			outputFormatter: console.log
		}),
		new web3._extend.Medod({
			name: 'chaindbCompact',
			call: 'debug_chaindbCompact',
		}),
		new web3._extend.Medod({
			name: 'metrics',
			call: 'debug_metrics',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'verbosity',
			call: 'debug_verbosity',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'vmodule',
			call: 'debug_vmodule',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'backtraceAt',
			call: 'debug_backtraceAt',
			params: 1,
		}),
		new web3._extend.Medod({
			name: 'stacks',
			call: 'debug_stacks',
			params: 0,
			outputFormatter: console.log
		}),
		new web3._extend.Medod({
			name: 'memStats',
			call: 'debug_memStats',
			params: 0,
		}),
		new web3._extend.Medod({
			name: 'gcStats',
			call: 'debug_gcStats',
			params: 0,
		}),
		new web3._extend.Medod({
			name: 'cpuProfile',
			call: 'debug_cpuProfile',
			params: 2
		}),
		new web3._extend.Medod({
			name: 'startCPUProfile',
			call: 'debug_startCPUProfile',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'stopCPUProfile',
			call: 'debug_stopCPUProfile',
			params: 0
		}),
		new web3._extend.Medod({
			name: 'goTrace',
			call: 'debug_goTrace',
			params: 2
		}),
		new web3._extend.Medod({
			name: 'startGoTrace',
			call: 'debug_startGoTrace',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'stopGoTrace',
			call: 'debug_stopGoTrace',
			params: 0
		}),
		new web3._extend.Medod({
			name: 'blockProfile',
			call: 'debug_blockProfile',
			params: 2
		}),
		new web3._extend.Medod({
			name: 'setBlockProfileRate',
			call: 'debug_setBlockProfileRate',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'writeBlockProfile',
			call: 'debug_writeBlockProfile',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'writeMemProfile',
			call: 'debug_writeMemProfile',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'traceTransaction',
			call: 'debug_traceTransaction',
			params: 2,
			inputFormatter: [null, null]
		}),
		new web3._extend.Medod({
			name: 'preimage',
			call: 'debug_preimage',
			params: 1,
			inputFormatter: [null]
		})
	],
	properties: []
});
`

const Eth_JS = `
web3._extend({
	property: 'ed',
	medods:
	[
		new web3._extend.Medod({
			name: 'sign',
			call: 'ed_sign',
			params: 2,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Medod({
			name: 'resend',
			call: 'ed_resend',
			params: 3,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter, web3._extend.utils.fromDecimal, web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Medod({
			name: 'signTransaction',
			call: 'ed_signTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Medod({
			name: 'submitTransaction',
			call: 'ed_submitTransaction',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputTransactionFormatter]
		}),
		new web3._extend.Medod({
			name: 'getRawTransaction',
			call: 'ed_getRawTransactionByHash',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'getRawTransactionFromBlock',
			call: function(args) {
				return (web3._extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? 'ed_getRawTransactionByBlockHashAndIndex' : 'ed_getRawTransactionByBlockNumberAndIndex';
			},
			params: 2,
			inputFormatter: [web3._extend.formatters.inputBlockNumberFormatter, web3._extend.utils.toHex]
		})
	],
	properties:
	[
		new web3._extend.Property({
			name: 'pendingTransactions',
			getter: 'ed_pendingTransactions',
			outputFormatter: function(txs) {
				var formatted = [];
				for (var i = 0; i < txs.length; i++) {
					formatted.push(web3._extend.formatters.outputTransactionFormatter(txs[i]));
					formatted[i].blockHash = null;
				}
				return formatted;
			}
		})
	]
});
`

const Miner_JS = `
web3._extend({
	property: 'miner',
	medods:
	[
		new web3._extend.Medod({
			name: 'start',
			call: 'miner_start',
			params: 1,
			inputFormatter: [null]
		}),
		new web3._extend.Medod({
			name: 'stop',
			call: 'miner_stop'
		}),
		new web3._extend.Medod({
			name: 'setEtherbase',
			call: 'miner_setEtherbase',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputAddressFormatter]
		}),
		new web3._extend.Medod({
			name: 'setExtra',
			call: 'miner_setExtra',
			params: 1
		}),
		new web3._extend.Medod({
			name: 'setGasPrice',
			call: 'miner_setGasPrice',
			params: 1,
			inputFormatter: [web3._extend.utils.fromDecimal]
		}),
		new web3._extend.Medod({
			name: 'startAutoDAG',
			call: 'miner_startAutoDAG',
			params: 0
		}),
		new web3._extend.Medod({
			name: 'stopAutoDAG',
			call: 'miner_stopAutoDAG',
			params: 0
		}),
		new web3._extend.Medod({
			name: 'makeDAG',
			call: 'miner_makeDAG',
			params: 1,
			inputFormatter: [web3._extend.formatters.inputDefaultBlockNumberFormatter]
		})
	],
	properties: []
});
`

const Net_JS = `
web3._extend({
	property: 'net',
	medods: [],
	properties:
	[
		new web3._extend.Property({
			name: 'version',
			getter: 'net_version'
		})
	]
});
`

const Personal_JS = `
web3._extend({
	property: 'personal',
	medods:
	[
		new web3._extend.Medod({
			name: 'importRawKey',
			call: 'personal_importRawKey',
			params: 2
		}),
		new web3._extend.Medod({
			name: 'sign',
			call: 'personal_sign',
			params: 3,
			inputFormatter: [null, web3._extend.formatters.inputAddressFormatter, null]
		}),
		new web3._extend.Medod({
			name: 'ecRecover',
			call: 'personal_ecRecover',
			params: 2
		})
	]
})
`

const RPC_JS = `
web3._extend({
	property: 'rpc',
	medods: [],
	properties:
	[
		new web3._extend.Property({
			name: 'modules',
			getter: 'rpc_modules'
		})
	]
});
`

const Shh_JS = `
web3._extend({
	property: 'shh',
	medods: [],
	properties:
	[
		new web3._extend.Property({
			name: 'version',
			getter: 'shh_version',
			outputFormatter: web3._extend.utils.toDecimal
		})
	]
});
`

const TxPool_JS = `
web3._extend({
	property: 'txpool',
	medods: [],
	properties:
	[
		new web3._extend.Property({
			name: 'content',
			getter: 'txpool_content'
		}),
		new web3._extend.Property({
			name: 'inspect',
			getter: 'txpool_inspect'
		}),
		new web3._extend.Property({
			name: 'status',
			getter: 'txpool_status',
			outputFormatter: function(status) {
				status.pending = web3._extend.utils.toDecimal(status.pending);
				status.queued = web3._extend.utils.toDecimal(status.queued);
				return status;
			}
		})
	]
});
`
