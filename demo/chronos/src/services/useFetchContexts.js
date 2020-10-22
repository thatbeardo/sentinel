import { useState, useEffect } from "react";

const resources = {
  jeff: `10833235-0135-4ee0-bc47-23820b5c52cb`,
  andy: `f8f40e31-5d0c-4276-8445-9ac8c4bd6ea2`,
  tasks: `9a4eb870-3c9b-46d2-8fee-17e55e335217`,
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
              "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJUVTJNRVZDTXpsQk9EQkJNa1U1TVVNeFEwWXdSakF4TlRNMFFqTXlOakJDUmpFd1JqZ3hOZyJ9.eyJpc3MiOiJodHRwczovL2JpdGhpcHBpZS5hdXRoMC5jb20vIiwic3ViIjoiZzQ3ZXhZcUR2emZJSlN1MEp3QUFYSlhuNVFqVFdxYkxAY2xpZW50cyIsImF1ZCI6Imh0dHBzOi8vYXBpLmd1YXJkbXkuYXBwLyIsImlhdCI6MTYwMzMwNTcxNiwiZXhwIjoxNjAzMzkyMTE2LCJhenAiOiJnNDdleFlxRHZ6ZklKU3UwSndBQVhKWG41UWpUV3FiTCIsImd0eSI6ImNsaWVudC1jcmVkZW50aWFscyJ9.Lz4lQsT7oppGJ-myG_aTFaZESNwzjp4PaBWjb25k_oJ0Xk7eGB0-W7MyQ3kJ3KPCqFTeotuiI-lLQSj_Hs366JrU9NRrvckUe12ANGKecKkbshUkrFvNq-O4GSCrHsiPdB5JQOBYpF8RqrrozWIWcajlXIIiuzRSpdH-z-CZcbYbIr6-s2lZ1KYggGp_snPzPvboddkW-gYj9aquDktXvvUBDdJ9LecB2uukkancUCqe3Az41l5SiIaNefnF9QxcEgdgpC2x-515FHK_nRTrw9j6qTL4Nt_2KeKdRr8f9TVxAeHD1SBQ92lY8iNKs4-2Cxprh3JdS91kXK7DEVVB9g",
          },
        }
      );
      const response = await apiResponse.json();
      if (response.data !== undefined) {
        setContexts(
          response.data.map((entitlements) => entitlements.attributes.name)
        );
      }
    };

    fetchData();
  }, [user]);

  return { contexts };
};
