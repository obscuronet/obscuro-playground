import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";
import {
  Batch,
  BatchDetails,
  BatchResponse,
} from "@/src/types/interfaces/BatchInterfaces";

export const fetchBatches = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<BatchResponse>> => {
  return await httpRequest<ResponseDataInterface<BatchResponse>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatches),
    searchParams: payload,
  });
};

export const fetchLatestBatch = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<Batch>> => {
  return await httpRequest<ResponseDataInterface<Batch>>({
    method: "get",
    url: pathToUrl(apiRoutes.getLatestBatch),
    searchParams: payload,
  });
};

export const fetchBatchByHash = async (
  hash: string
): Promise<ResponseDataInterface<BatchDetails>> => {
  return await httpRequest<ResponseDataInterface<BatchDetails>>({
    method: "get",
    url: pathToUrl(apiRoutes.getBatchByHash, { hash }),
  });
};
