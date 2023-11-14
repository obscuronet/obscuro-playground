import { httpRequest } from ".";
import { apiRoutes } from "@/src/routes";
import { pathToUrl } from "@/src/routes/router";
import { ResponseDataInterface } from "@/src/types/interfaces";

export const getRollups = async (
  payload?: Record<string, any>
): Promise<ResponseDataInterface<any>> => {
  const data = await httpRequest<ResponseDataInterface<any>>({
    method: "get",
    url: pathToUrl(apiRoutes.getRollups),
    searchParams: payload,
  });
  return data;
};