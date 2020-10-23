import { useState, useEffect } from "react";

const resources = {
  jeff: `10833235-0135-4ee0-bc47-23820b5c52cb`,
  andy: `f8f40e31-5d0c-4276-8445-9ac8c4bd6ea2`,
  tasks: `9a4eb870-3c9b-46d2-8fee-17e55e335217`,
};

const defaultContexts = {
  jeff: `af40bb8a-4343-428d-a1e2-728cad3668cf`,
  andy: `5ff83647-3d9e-46b8-b5ee-e798f76ef5db`,
};

export const useFetchContexts = (user) => {
  const [contexts, setContexts] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      const apiResponse = await fetch(
        `http://localhost:8081/v1/resources/${resources[user]}/contexts?claimant=amazon-dev-ops`,
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
        setContexts(
          response.data.map((entitlements) => {
            let value = entitlements.attributes.name;
            let name =
              entitlements.id === defaultContexts[user]
                ? `${value} (default)`
                : value;
            return { name, value };
          })
        );
      }
    };

    fetchData();
  }, [user]);

  return { contexts };
};
