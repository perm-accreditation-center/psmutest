export interface Question {
  id: number;
  question: string;
  options: string[];
  correctAnswer: number;
}

export interface Test {
  id: number;
  title: string;
  questions: Question[];
}

export interface TestResult {
  userId: string;
  firstName: string;
  lastName: string;
  middleName?: string;
  testId: string | number;
  answers: Record<number, number>;
  score?: number;
  date?: string;
}

export interface UserData {
  firstName: string;
  lastName: string;
  middleName: string;
  userId: string;
}
