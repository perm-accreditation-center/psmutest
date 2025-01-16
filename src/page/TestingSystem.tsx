import React from "react";
import { Box, Alert, CircularProgress, Button } from "@mui/material";
import { FinalPage } from "./FinalStep";
import { ProgressStepper } from "../components/ProgressStepper";
import { TestSection } from "../components/TestSection";
import { UserDataForm } from "../components/UserDataForm";
import { ConfirmDialog } from "../components/ConfirmDialog";
import { useTestingSystem } from "./useTestingSystem";

export const TestingSystem: React.FC = () => {
  const {
    loading,
    initialized,
    activeStep,
    setResetDialogOpen,
    error,
    tests,
    userData,
    setUserData,
    setActiveStep,
    testAnswers,
    setTestAnswers,
    handleTestSubmit,
    resetDialogOpen,
    handleReset,
  } = useTestingSystem();

  if (loading && !initialized) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="100dvh"
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 1200, mx: "auto", p: { xs: 2, sm: 4 } }}>
      <Box sx={{ display: "flex", justifyContent: "flex-end", mb: 2 }}>
        <Button
          color="error"
          variant="outlined"
          disabled={activeStep === 0}
          onClick={() => setResetDialogOpen(true)}
        >
          Завершить тестирование
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <ProgressStepper tests={tests} activeStep={activeStep} />

      {activeStep === 0 ? (
        <UserDataForm
          userData={userData}
          onChange={setUserData}
          onSubmit={() => setActiveStep(1)}
        />
      ) : activeStep <= tests.length ? (
        <TestSection
          test={tests[activeStep - 1]}
          answers={testAnswers[tests[activeStep - 1].id] || {}}
          onAnswerChange={(questionId, answer) => {
            setTestAnswers((prev) => ({
              ...prev,
              [tests[activeStep - 1].id]: {
                ...(prev[tests[activeStep - 1].id] || {}),
                [questionId]: answer,
              },
            }));
          }}
          onSubmit={handleTestSubmit}
        />
      ) : (
        <FinalPage userId={userData.userId} tests={tests} />
      )}

      <ConfirmDialog
        open={resetDialogOpen}
        onClose={() => setResetDialogOpen(false)}
        onConfirm={handleReset}
        title="Завершить тестирование?"
        message="Все несохраненные результаты будут потеряны. Вы уверены, что хотите завершить тестирование?"
      />
    </Box>
  );
};
