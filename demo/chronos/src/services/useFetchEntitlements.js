import { useState, useEffect } from "react";

const resources = {
  jeff: `10833235-0135-4ee0-bc47-23820b5c52cb`,
  andy: `f8f40e31-5d0c-4276-8445-9ac8c4bd6ea2`,
  tasks: `9a4eb870-3c9b-46d2-8fee-17e55e335217`,
};

const salesRegions = {
  // North America
  "3fe62f7f-e3cc-4830-b0dd-e19347d48ffe": "USA",
  "fe3a7f6c-9795-477a-8efa-569dd53b3c4a": "CU",
  "c8ecd967-00c1-484d-8618-71e755759636": "CAN",
  "1d7fcf8e-2b69-4888-ab5c-0f837a9902b3": "PA",

  // Asia Pacific
  "f24f7485-783f-42f0-8378-54a084c82ab2": "IND",
  "a5be3007-49d6-43e2-aa7d-b3f43f923fc5": "TAI",
  "a7ba149e-59c9-4a0c-9c5a-8ff15227d710": "NZ",
  "06dd3619-fd6c-4dc1-8ba6-28db4699a35e": "SG",
};

const contexts = {
  ceo: `af40bb8a-4343-428d-a1e2-728cad3668cf`,
  "vp-apac": `5ff83647-3d9e-46b8-b5ee-e798f76ef5db`,
};

export const useFetchEntitlements = (user, context) => {
  context = context === "" ? "" : contexts[context];
  const [permitted, setPermitted] = useState(false);
  const [regions, setRegions] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      const apiResponse = await fetch(
        `http://localhost:8081/v1/principal/${resources[user]}/authorization?claimant=amazon-dev-ops&permissions=read&context_id=${context}&depth=5&include_denied=false`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            "x-sentinel-tenant": "dev",
            Authorization:
              "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJUVTJNRVZDTXpsQk9EQkJNa1U1TVVNeFEwWXdSakF4TlRNMFFqTXlOakJDUmpFd1JqZ3hOZyJ9.eyJpc3MiOiJodHRwczovL2JpdGhpcHBpZS5hdXRoMC5jb20vIiwic3ViIjoiZzQ3ZXhZcUR2emZJSlN1MEp3QUFYSlhuNVFqVFdxYkxAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vYXBpLmd1YXJkbXkuYXBwLyIsImlhdCI6MTYwMzM3MDEyMiwiZXhwIjoxNjAzNDU2NTIyLCJhenAiOiJnNDdleFlxRHZ6ZklKU3UwSndBQVhKWG41UWpUV3FiTCIsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyJ9.z7EhonVEY5-bPHa07QXT9_MpLg2QN5U0LTk2RUWLK2Ewh5xP2BWQ5pJZXoQ2CXri562zUNd3wmmLg-fdg2mRwZg1i5aTrOFJVrtcaQWrRHUTxLkHcGL5MEvDNX9yMyAaehqSupGaLGW68Q-66jAWSGjfL39aIe42yHLDpkc97B-4atuHtqgp0FJHAtnxPk9NqgGZ-Dmx0TTGdHdu-5QX9x0kf2FVbV6p97XkGK7-JpKVKWq8PBzskw3fVH56gj05VAqeYXCA-WPQIIaivFcys_SqSsNYVTpjN6bmupQnZW_zQwEVSNkxjRY9fyzS3zWSMKojAn7xu7nBbgiUD7qomg",
          },
        }
      );
      const response = await apiResponse.json();

      if (response.data !== undefined) {
        setPermitted(
          response.data.filter(
            (entitlement) => entitlement.id === resources["tasks"]
          ).length > 0
        );
        setRegions(
          response.data
            .filter((entitlement) => entitlement.id in salesRegions)
            .map((entitlement) => salesRegions[entitlement.id])
        );
      }
    };

    fetchData();
  }, [user, context]);

  return { permitted, regions };
};
