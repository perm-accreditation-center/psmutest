import React from "react";
import { Box, TextField, Button } from "@mui/material";
import { UserData } from "../types";

interface UserDataFormProps {
  userData: UserData;
  onChange: (data: UserData) => void;
  onSubmit: () => void;
}

export const UserDataForm: React.FC<UserDataFormProps> = ({
  userData,
  onChange,
  onSubmit,
}) => {
  const handleChange =
    (field: keyof UserData) => (e: React.ChangeEvent<HTMLInputElement>) => {
      onChange({ ...userData, [field]: e.target.value });
    };

  const isFormValid = userData.firstName.trim() && userData.lastName.trim();

  return (
    <Box
      component="form"
      display="flex"
      flexDirection="column"
      gap={3}
      maxWidth={400}
      mx="auto"
      p={3}
    >
      <TextField
        label="Имя"
        value={userData.firstName}
        onChange={handleChange("firstName")}
        required
        autoComplete="off"
        autoFocus={true}
        fullWidth
        variant="outlined"
      />
      <TextField
        label="Фамилия"
        value={userData.lastName}
        onChange={handleChange("lastName")}
        required
        autoComplete="off"
        fullWidth
        variant="outlined"
      />
      <TextField
        label="Отчество"
        value={userData.middleName}
        onChange={handleChange("middleName")}
        fullWidth
        autoComplete="off"
        variant="outlined"
      />
      <Button
        variant="contained"
        onClick={onSubmit}
        disabled={!isFormValid}
        size="large"
      >
        Продолжить
      </Button>
    </Box>
  );
};
