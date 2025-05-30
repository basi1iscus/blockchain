// Utilities
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    wallet: null as null | {
      address: string
      public_key: string
      public_key_hash: string
      private_key: string
    },
    scriptResult: null as null | { code: string, result: boolean, success: boolean },
    compileResult: null as null | { scriptSig: string, scriptPubKey: string, success: boolean },
    parseResult: null as null | { scriptSig: string, scriptPubKey: string, success: boolean },
    loading: false,
    error: null as null | string,
  }),
  actions: {
    async createWallet () {
      this.loading = true
      this.error = null
      try {
        const res = await fetch('/api/wallet', { method: 'POST' })
        if (!res.ok) throw new Error('Failed to create wallet')
        const data = await res.json()
        this.wallet = data
      } catch (e: any) {
        this.error = e.message || 'Unknown error'
      } finally {
        this.loading = false
      }
    },
    async runScript (payload: { scriptSig: string, scriptPubKey: string, signedData: string }) {
      this.loading = true
      this.error = null
      try {
        const res = await fetch(`/api/sript/run`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            script_sig: payload.scriptSig,
            script_pub_key: payload.scriptPubKey,
            signed_data: payload.signedData,
          }),
        })
        if (!res.ok) throw new Error('Failed to run script')
        const data = await res.json()
        this.scriptResult = data
      } catch (e: any) {
        this.error = e.message || 'Unknown error'
      } finally {
        this.loading = false
      }
    },
    async compileScript (payload: { scriptSig: string, scriptPubKey: string}) {
      this.loading = true
      this.error = null
      try {
        const res = await fetch(`/api/sript/compile`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            script_sig: payload.scriptSig,
            script_pub_key: payload.scriptPubKey,
          }),
        })
        if (!res.ok) throw new Error('Failed to run script')
        const data = await res.json()
        this.compileResult = data
      } catch (e: any) {
        this.error = e.message || 'Unknown error'
      } finally {
        this.loading = false
      }
    },
    async parseScript (payload: { scriptSig: string, scriptPubKey: string}) {
      this.loading = true
      this.error = null
      try {
        const res = await fetch(`/api/sript/parse`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            script_sig: payload.scriptSig,
            script_pub_key: payload.scriptPubKey,
          }),
        })
        if (!res.ok) throw new Error('Failed to run script')
        const data = await res.json()
        this.parseResult = data
      } catch (e: any) {
        this.error = e.message || 'Unknown error'
      } finally {
        this.loading = false
      }
    },
  },
})
