import React from "react";
import { Box, Button, Typography } from "@mui/material";
import { Test } from "../types";
import { QuestionCard } from "./QuestionCard";

interface TestSectionProps {
  test: Test;
  answers: Record<number, number>;
  onAnswerChange: (questionId: number, answer: number) => void;
  onSubmit: () => void;
}

export const TestSection: React.FC<TestSectionProps> = ({
  test,
  answers,
  onAnswerChange,
  onSubmit,
}) => {
  const isTestComplete = test.questions.every(
    (q) => answers[q.id] !== undefined
  );

  return (
    <Box maxWidth={800} mx="auto" p={3}>
      <Typography variant="h4" marginBottom={5} fontWeight={600} align="center">
        {test.title}
      </Typography>
      {test.questions.map((question) => (
        <QuestionCard
          key={question.id}
          question={question}
          selectedAnswer={answers[question.id]}
          onAnswerSelect={(answer) => onAnswerChange(question.id, answer)}
        />
      ))}
      <Box display="flex" justifyContent="center" mt={4}>
        <Button
          variant="contained"
          onClick={() => onSubmit()}
          disabled={!isTestComplete}
          size="large"
        >
          Завершить тест
        </Button>
      </Box>
    </Box>
  );
};
