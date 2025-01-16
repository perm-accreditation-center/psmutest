import React from "react";
import { Box, Button, Paper, Typography } from "@mui/material";
import { Question } from "../types";

interface QuestionCardProps {
  question: Question;
  selectedAnswer?: number;
  onAnswerSelect: (answer: number) => void;
}

export const QuestionCard: React.FC<QuestionCardProps> = ({
  question,
  selectedAnswer,
  onAnswerSelect,
}) => {
  return (
    <Paper elevation={2} sx={{ p: 3, mb: 3 }}>
      <Typography variant="h6" gutterBottom>
        {question.question}
      </Typography>
      <Box display="flex" flexDirection="column" gap={1} mt={2}>
        {question.options.map((option, idx) => (
          <Button
            key={idx}
            variant={selectedAnswer === idx ? "contained" : "outlined"}
            onClick={() => onAnswerSelect(idx)}
            fullWidth
            sx={{ justifyContent: "flex-start", px: 2, py: 1.5 }}
          >
            {option}
          </Button>
        ))}
      </Box>
    </Paper>
  );
};
