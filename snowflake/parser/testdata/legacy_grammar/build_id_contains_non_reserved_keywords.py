""" build_id_contains_non_reserved_keywords.py rebuilds the ID and ID2 token to contains the non-reserved keywords.

It does the following steps in sequence:
1.  extracts the token in SnowflakeLexer.g4 from "// Build id contains the build non reserved keywords start." 
    to "// Build id contains The token for "the build non reserved keywords stop" is extracted.
2.  At the end of SnowflakeParser.g4, write the supplement_non_reserved_words rule, 
    and add supplement_non_reserved_words to the rule for id_ as a candidate.

Usage:
python3 build_id_contains_non_reserved_keywords.py
"""

import re

snowflake_reserved_keyword = {
    "ACCOUNT": True,
    "ALL": True,
    "ALTER": True,
    "AND": True,
    "ANY": True,
    "AS": True,

    "BETWEEN": True,
    "BY": True,

    "CASE": True,
    "CAST": True,
    "CHECK": True,
    "COLUMN": True,
    "CONNECT": True,
    "CONNECTION": True,
    "CONSTRAINT": True,
    "CREATE": True,
    "CROSS": True,
    "CURRENT": True,
    "CURRENT_DATE": True,
    "CURRENT_TIME": True,
    "CURRENT_TIMESTAMP": True,
    "CURRENT_USER": True,

    "DATABASE": True,
    "DELETE": True,
    "DISTINCT": True,
    "DROP": True,

    "ELSE": True,
    "EXISTS": True,
    "FALSE": True,
    "FOLLOWING": True,
    "FOR": True,
    "FROM": True,
    "FULL": True,

    "GRANT": True,
    "GROUP": True,
    "GSCLUSTER": True,

    "HAVING": True,

    "ILIKE": True,
    "IN": True,
    "INCREMENT": True,
    "INNER": True,
    "INSERT": True,
    "INTERSECT": True,
    "INTO": True,
    "IS": True,
    "ISSUE": True,

    "JOIN": True,

    "LATERAL": True,
    "LEFT": True,
    "LIKE": True,
    "LOCALTIME": True,
    "LOCALTIMESTAMP": True,

    "MINUS": True,

    "NATURAL": True,
    "NOT": True,
    "NULL": True,

    "OF": True,
    "ON": True,
    "OR": True,
    "ORDER": True,
    "ORGANIZATION": True,

    "QUALIFY": True,

    "REGEXP": True,
    "REVOKE": True,
    "RIGHT": True,
    "RLIKE": True,
    "ROW": True,
    "ROWS": True,

    "SAMPLE": True,
    "SCHEMA": True,
    "SELECT": True,
    "SET": True,
    "SOME": True,
    "START": True,

    "TABLE": True,
    "TABLESAMPLE": True,
    "THEN": True,
    "TO": True,
    "TRIGGER": True,
    "TRUE": True,
    "TRY_CAST": True,

    "UNION": True,
    "UNIQUE": True,
    "UPDATE": True,
    "USING": True,

    "VALUES": True,
    "VIEW": True,

    "WHEN": True,
    "WHENEVER": True,
    "WHERE": True,
    "WITH": True
}

def read_tokens_name_before_token_from_lexer_file(filepath: str, token: str) -> list[str]:    
    tokens_name_before_token = []
    regex = r"^\s*(?P<token_name>[A-Z_][A-Z_0-9]*)\s*:"
    start_placeholder = "Build id contains the non reserved keywords start."
    stop_placeholder = "Build id contains the non reserved keywords stop."
    begin = False
    with open(filepath, "r") as lexer_file:
        lines = lexer_file.readlines()
        has_placeholder = any(start_placeholder in line for line in lines) and any(stop_placeholder in line for line in lines)
        if not has_placeholder:
            begin = True
        reached_target_token = False
        for line in lines:
            if has_placeholder and start_placeholder in line:
                begin = True
                continue
            if line.isspace() or (not begin):
                continue
            if has_placeholder and (stop_placeholder in line):
                break

            matches = re.finditer(regex, line, re.MULTILINE)
            for matchNum, match in enumerate(matches, start=1):
                if matchNum > 1:
                    break
                if match.group("token_name") == token:
                    reached_target_token = True
                    break
                tokens_name_before_token.append(match.group("token_name"))
            if reached_target_token:
                break
    return tokens_name_before_token

def pretty_print(tokens: list[str], hello: str | None) -> None:
    if hello is not None:
        print(hello)
    # Format:
    # item: [N]: token_name
    # 5 items in one line, N is the index of the token in the list.
    for index, element in enumerate(tokens):
        print(f"[{index}]: {element}", end=", " if (index+1) % 5 != 0 and index != len(tokens) - 1 else "\n")

    print()

def append_non_reserved_token_to_rules_in_parser(parser_file_path: str, append_rules_token_name: str, new_token_name: str, token_list: list[str]):
    content = ""
    with open(parser_file_path, "r+") as file:
        content = file.read()
    token_content = get_content_by_token_name(content, append_rules_token_name)
    if token_content is None:
        print(f"Cannot find token {append_rules_token_name} in {parser_file_path}")
        return
    # Remove the last semicolon.
    new_token_content = token_content[:-1]
    # Append new token name.
    new_token_content += f"| {new_token_name}\n    ;"
    if new_token_name not in token_content:
        content = content.replace(token_content, new_token_content)

    # Append new token rule.
    new_token_rule = ""
    for index, identifier in enumerate(token_list):
        if index == 0:
            new_token_rule += f"{new_token_name}\n    : {identifier}"
        else:
            new_token_rule += f"\n    | {identifier}"
    new_token_rule += f"\n    ;"
    already_exist_new_token_rule = get_content_by_token_name(content, new_token_name)
    if already_exist_new_token_rule:
        content = content.replace(already_exist_new_token_rule, new_token_rule)
    else:
        content = content + "\n" + new_token_rule + "\n"
    with open(parser_file_path, "w") as file:
        file.write(content)
        file.flush()
    


def get_content_by_token_name(content: str, token_name: str) -> str: 
    # Support upstream grammar formatting where comments may appear between
    # rule name and ':' and where rules can be indented.
    token_regex = r"(?ms)^\s*(%s)\s*(?:\n\s*//[^\n]*)*\n\s*:[\s\S]*?\n\s*;" % re.escape(token_name)
    token_content = re.search(token_regex, content)
    if token_content:
        return token_content.group(0)
    return None

        
    

if __name__ == "__main__":
    tokens = read_tokens_name_before_token_from_lexer_file("SnowflakeLexer.g4", "ID")
    pretty_print(tokens, "Tokens before ID token:")
    filtered_tokens = [item for item in tokens if item not in snowflake_reserved_keyword]
    pretty_print(filtered_tokens, "Tokens before ID token without reserved keywords:")
    append_non_reserved_token_to_rules_in_parser("SnowflakeParser.g4", "id_", "supplement_non_reserved_words", filtered_tokens)
    
