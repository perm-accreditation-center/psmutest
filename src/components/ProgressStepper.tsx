import React from "react";
import { Stepper, Step, StepLabel } from "@mui/material";
import { Test } from "../types";

interface ProgressStepperProps {
  tests: Test[];
  activeStep: number;
}

export const ProgressStepper: React.FC<ProgressStepperProps> = ({
  tests,
  activeStep,
}) => {
  const steps = [
    "Персональные данные",
    ...tests.map((test) => `${test.title}`),
    "Завершение",
  ];

  return (
    <Stepper
      activeStep={activeStep}
      alternativeLabel
      sx={{ mb: 4, px: { xs: 1, sm: 4 } }}
    >
      {steps.map((label) => (
        <Step key={label}>
          <StepLabel>{label}</StepLabel>
        </Step>
      ))}
    </Stepper>
  );
};
