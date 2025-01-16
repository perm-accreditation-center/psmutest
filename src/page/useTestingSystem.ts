import { useState, useEffect } from "react";
import { getTests, submitTestResult } from "../api/api";
import { Test, UserData, TestResult } from "../types";

interface StoredState {
  userData: UserData;
  testAnswers: Record<string | number, Record<number, number>>;
  activeStep: number;
  completedSteps: Record<number, boolean>;
}

export function useTestingSystem() {
  const [initialized, setInitialized] = useState(false);
  const [tests, setTests] = useState<Test[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [resetDialogOpen, setResetDialogOpen] = useState(false);

  const [userData, setUserData] = useState<UserData>({
    firstName: "",
    lastName: "",
    middleName: "",
    userId: crypto.randomUUID(),
  });
  const [testAnswers, setTestAnswers] = useState<
    Record<string | number, Record<number, number>>
  >({});
  const [activeStep, setActiveStep] = useState(0);
  const [completedSteps, setCompletedSteps] = useState<Record<number, boolean>>(
    {}
  );

  useEffect(() => {
    const loadInitialData = async () => {
      try {
        setLoading(true);

        const savedState = localStorage.getItem("testingState");
        if (savedState) {
          const parsed: StoredState = JSON.parse(savedState);
          setUserData(parsed.userData);
          setTestAnswers(parsed.testAnswers);
          setActiveStep(parsed.activeStep);
          setCompletedSteps(parsed.completedSteps);
        }

        const response = await getTests();
        setTests(response.data.tests);
        setInitialized(true);
      } catch (err) {
        console.error(err);
        setError("Ошибка загрузки данных");
      } finally {
        setLoading(false);
      }
    };

    loadInitialData();
  }, []);

  const handleReset = () => {
    setTimeout(() => {
      localStorage.removeItem("testingState");
      location.reload();
    }, 700);
  };

  useEffect(() => {
    if (!initialized) return;

    const state: StoredState = {
      userData,
      testAnswers,
      activeStep,
      completedSteps,
    };

    localStorage.setItem("testingState", JSON.stringify(state));
  }, [userData, testAnswers, activeStep, completedSteps, initialized]);

  const handleStepComplete = (stepIndex: number) => {
    setCompletedSteps((prev) => ({
      ...prev,
      [stepIndex]: true,
    }));
  };

  const handleTestSubmit = async () => {
    scrollToTopWithDelay(300); // Задержка 300 мс

    const currentTest = tests[activeStep - 1];
    const rawAnswers = testAnswers[currentTest.id] || {};

    const incrementedAnswers = Object.fromEntries(
      Object.entries(rawAnswers).map(([questionId, answer]) => [
        questionId,
        answer + 1,
      ])
    );

    const result: TestResult = {
      ...userData,
      testId: currentTest.id,
      answers: incrementedAnswers,
    };

    try {
      setLoading(true);
      await submitTestResult(result);
      handleStepComplete(activeStep);
      setActiveStep((prev) => prev + 1);
    } catch (err) {
      console.error(err);
      setError("Ошибка при отправке результатов для текущего теста");
    } finally {
      setLoading(false);
    }
  };

  const scrollToTopWithDelay = (delay = 0) => {
    setTimeout(() => {
      if ("scrollBehavior" in document.documentElement.style) {
        window.scrollTo({ top: 0, behavior: "smooth" });
      } else {
        window.scrollTo(0, 0);
      }
    }, delay);
  };

  const handleFinalSubmit = async () => {
    try {
      setLoading(true);

      handleReset();
    } catch (err) {
      console.error(err);
      setError("Ошибка при отправке результатов");
    } finally {
      setLoading(false);
    }
  };

  return {
    tests,
    loading,
    initialized,
    error,
    resetDialogOpen,
    setResetDialogOpen,
    handleReset,
    userData,
    setUserData,
    testAnswers,
    setTestAnswers,
    activeStep,
    setActiveStep,
    completedSteps,
    handleStepComplete,
    handleTestSubmit,
    handleFinalSubmit,
  };
}
