module Main where

newtype Parser a = Parser (String -> [(a, String)])

item :: Parser Char
item = Parser (\cs -> case cs of 
  ""      -> []
  (c:cs)  -> [(c, cs)])

instance Monad Parser where
  return a  = Parser (\cs -> [(a, cs)])
  p >>= f   = Parser (\cs -> concat [parse (f a) cs' | (a,cs') <- parse p cs])