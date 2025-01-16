import React, { useEffect, useState } from "react";
import { Box, Typography, CircularProgress } from "@mui/material";

interface TestResultProps {
  title: string;
  percentage: number; // Процент результата
}

export const TestResult: React.FC<TestResultProps> = ({ title, percentage }) => {
  const [animatedPercentage, setAnimatedPercentage] = useState(0);

  useEffect(() => {
    const step = percentage / 100;
    let progress = 0;

    const interval = setInterval(() => {
      progress += step;
      if (progress >= percentage) {
        setAnimatedPercentage(percentage);
        clearInterval(interval);
      } else {
        setAnimatedPercentage(Math.floor(progress));
      }
    }, 10);

    return () => clearInterval(interval);
  }, [percentage]);

  return (
    <Box textAlign="center" mb={4}>
      <Typography variant="h6" gutterBottom>
        {title}
      </Typography>
      <Box position="relative" display="inline-flex">
        <CircularProgress
          variant="determinate"
          value={animatedPercentage}
          size={120}
          thickness={5}
        />
        <Box
          top={0}
          left={0}
          bottom={0}
          right={0}
          position="absolute"
          display="flex"
          alignItems="center"
          justifyContent="center"
        >
          <Typography variant="h4" component="div" color="textPrimary">
            {animatedPercentage}%
          </Typography>
        </Box>
      </Box>
    </Box>
  );
};
