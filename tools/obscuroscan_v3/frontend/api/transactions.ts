import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";
import {
  TransactionCount,
  Price,
  TransactionResponse,
} from "@/src/types/interfaces/TransactionInterfaces";

export const getTransactions = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<TransactionResponse>> => {
  const data = await httpRequest<ResponseDataInterface<TransactionResponse>>({
    method: "get",
    url: pathToUrl(apiRoutes.getTransactions),
    searchParams: payload,
  });
  return data;
};

export const getTransactionCount = async (): Promise<TransactionCount> => {
  const data = await httpRequest<TransactionCount>({
    method: "get",
    url: pathToUrl(apiRoutes.getTransactionCount),
  });
  return data;
};

export const getPrice = async (): Promise<Price> => {
  const data = await httpRequest<Price>({
    method: "get",
    url: apiRoutes.getPrice,
  });
  return data;
};