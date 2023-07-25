import { defineStore } from 'pinia';
import Config from "@/lib/config";
import BatchList from "@/lib/batchList";

export const useBatchStore = defineStore({
    id: 'batchStore',
    state: () => ({
        latestBatch: null,
        latestL1Proof: null,
        batches: new BatchList(),
        loading: false,
        pollingInterval: 1000,  // 5 seconds
        timer: null,
    }),
    actions: {
        async fetchCount() {
            this.loading = true;
            try {
                let response = await fetch( Config.backendServerAddress+'/items/batch/latest/');
                let data = await response.json();
                this.latestBatch = data.item.Number;
                this.latestL1Proof = data.item.L1Proof;

                this.batches.add(data.item);

                console.log("Fetched "+this.latestBatch);
            } catch (error) {
                console.error("Failed to fetch count:", error);
            } finally {
                this.loading = false;
            }
        },

        startPolling() {
            this.stopPolling(); // Ensure previous intervals are cleared
            this.timer = setInterval(async () => {
                await this.fetchCount();
            }, this.pollingInterval);
        },

        stopPolling() {
            if (this.timer) {
                clearInterval(this.timer);
                this.timer = null;
            }
        }

    },
});
