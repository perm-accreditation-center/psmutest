import React, { useEffect, useState } from "react";
import { Box, Typography, CircularProgress, Alert } from "@mui/material";
import { getTestResults } from "../api/api"; // Запрос результатов
import { TestResult } from "../components/TestResult";
import { Test } from "../types";

interface FinalPageProps {
  userId: string;
  tests:  Test[]
}


export const FinalPage: React.FC<FinalPageProps> = ({ userId, tests }) => {
  const [results, setResults] = useState<
    { testId: number; title: string; percentage: number }[]
  >([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchResults = async () => {
      try {
        const fetchedResults = await Promise.all(
          tests.map(async (test) => {
            const response = await getTestResults(userId, test.id);
            return {
              testId: test.id,
              title: test.title,
              percentage: response.data.percentage,
            };
          })
        );
        setResults(fetchedResults);
      } catch (err) {
        console.log(err);
        setError("Ошибка загрузки результатов.");
      } finally {
        setLoading(false);
      }
    };

    fetchResults();
  }, [userId, tests]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box textAlign="center" p={4}>
      <Typography variant="h4" mb={4}>
        Ваши результаты
      </Typography>
      {results.map((result) => (
        <TestResult
          key={result.testId}
          title={result.title}
          percentage={result.percentage}
        />
      ))}
    </Box>
  );
};
