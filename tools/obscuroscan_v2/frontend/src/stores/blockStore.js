import { defineStore } from 'pinia';
import Config from "@/lib/config";
import BlockList from "@/lib/blockList";

export const useBlockStore = defineStore({
    id: 'blockStore',
    state: () => ({
        blocks: new BlockList(),
        loading: false,
        pollingInterval: Config.pollingInterval,
        timer: null,
    }),
    actions: {
        async fetchCount() {
            this.loading = true;
            try {
                let response = await fetch( Config.backendServerAddress+'/items/block/latest/');
                let data = await response.json();
                this.blocks.add(data.item);

                console.log("Fetched "+data.item.Number);
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
