import TruncatedAddress from "../common/truncated-address";
import { Avatar, AvatarFallback } from "@/src/components/ui/avatar";
import { Transaction } from "@/src/types/interfaces/TransactionInterfaces";
import { Badge } from "../../ui/badge";

export function RecentTransactions({ transactions }: { transactions: any }) {
  return (
    <div className="space-y-8">
      {transactions?.result?.TransactionsData.map(
        (transaction: Transaction, i: number) => (
          <div className="flex items-center" key={i}>
            <Avatar className="h-9 w-9">
              <AvatarFallback>TX</AvatarFallback>
            </Avatar>
            <div className="ml-4 space-y-1">
              <p className="text-sm font-medium leading-none">
                #{transaction?.BatchHeight}
              </p>
            </div>
            <div className="ml-auto font-medium">
              <TruncatedAddress address={transaction?.TransactionHash} />
            </div>
            <div className="ml-auto">
              <Badge>{transaction?.Finality}</Badge>
            </div>
          </div>
        )
      )}
    </div>
  );
}
