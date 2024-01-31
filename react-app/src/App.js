import React, { useState } from "react";
import {
  Container,
  TextField,
  Button,
  Typography,
  Paper,
  CircularProgress,
} from "@mui/material";
import { styled } from "@mui/system";

const StyledContainer = styled(Container)({
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  justifyContent: "center",
  minHeight: "100vh",
  background: "#333",
  color: "#fff",
  maxWidth: "100%!important",
});

const StyledPaper = styled(Paper)({
  padding: "20px",
  maxWidth: "600px",
  width: "100%",
  textAlign: "center",
  background: "#444",
});

const App = () => {
  const [inputText, setInputText] = useState("");
  const [translation, setTranslation] = useState("");
  const [loading, setLoading] = useState(false);

  const handleTranslate = async () => {
    try {
      setLoading(true);
      const response = await fetch("http://localhost:8080/translate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ text: inputText }),
      });
      const data = await response.json();
      setTranslation(data.translation);
    } catch (error) {
      console.error("Translation error:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <StyledContainer>
      <StyledPaper elevation={3}>
        <Typography variant="h4" gutterBottom style={{ color: "#fff" }}>
          MyPolyglot
        </Typography>
        <TextField
          label="Enter Text"
          variant="outlined"
          fullWidth
          multiline
          rows={4}
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          style={{ marginBottom: "16px" }}
          InputLabelProps={{
            style: { color: "#fff" },
          }}
          InputProps={{
            style: {
              color: "#fff",
              borderColor: "#888",
            },
          }}
        />
        <Button
          variant="contained"
          color="primary"
          onClick={handleTranslate}
          disabled={loading}
        >
          {loading ? (
            <CircularProgress size={24} style={{ color: "#fff" }} />
          ) : (
            "Translate"
          )}
        </Button>
        {translation && (
          <div style={{ marginTop: "16px" }}>
            <Typography variant="h6" gutterBottom style={{ color: "#fff" }}>
              Translation Result:
            </Typography>
            <Typography variant="body1" style={{ color: "#fff" }}>
              {translation}
            </Typography>
          </div>
        )}
      </StyledPaper>
    </StyledContainer>
  );
};

export default App;
