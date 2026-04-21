using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text.Json;
using Microsoft.SqlServer.TransactSql.ScriptDom;

namespace Bytebase.Omni.MssqlHarness;

// Harness: read SQL from stdin, parse with SqlScriptDOM, emit a compact
// JSON shape covering the fields we diff against omni's AST:
//   - select list element kinds + aliases
//   - FROM table-reference kinds + aliases + column-list contents
//
// Protocol: one SQL per invocation (stdin). Caller is responsible for
// starting/stopping the process per fixture. For batch use, prefer the
// line-delimited protocol below (env MSSQL_HARNESS_LINE=1 reads one SQL
// per line, emits one JSON per line until EOF).

public static class Program
{
    public static int Main()
    {
        var lineMode = Environment.GetEnvironmentVariable("MSSQL_HARNESS_LINE") == "1";
        return lineMode ? RunLineMode() : RunSingle();
    }

    private static int RunSingle()
    {
        var sql = Console.In.ReadToEnd();
        var json = ParseToJson(sql);
        Console.WriteLine(json);
        return 0;
    }

    private static int RunLineMode()
    {
        string? line;
        while ((line = Console.In.ReadLine()) != null)
        {
            // Line-mode encoding: base64 of the SQL, one per line.
            // This avoids newline-in-SQL ambiguity.
            try
            {
                var bytes = Convert.FromBase64String(line.Trim());
                var sql = System.Text.Encoding.UTF8.GetString(bytes);
                Console.WriteLine(ParseToJson(sql));
            }
            catch (FormatException)
            {
                Console.WriteLine("{\"error\":\"bad base64\"}");
            }
            Console.Out.Flush();
        }
        return 0;
    }

    private static string ParseToJson(string sql)
    {
        var parser = new TSql160Parser(initialQuotedIdentifiers: true);
        IList<ParseError> errors;
        var tree = parser.Parse(new StringReader(sql), out errors);
        var output = new Dictionary<string, object?>();
        if (errors.Count > 0)
        {
            output["errors"] = errors.Select(e => new Dictionary<string, object?>
            {
                ["line"] = e.Line,
                ["column"] = e.Column,
                ["message"] = e.Message,
            }).ToList();
        }
        output["shape"] = Shape.Of(tree);
        return JsonSerializer.Serialize(output, new JsonSerializerOptions
        {
            WriteIndented = false,
            DefaultIgnoreCondition = System.Text.Json.Serialization.JsonIgnoreCondition.WhenWritingNull,
        });
    }
}

internal static class Shape
{
    public static object? Of(TSqlFragment? frag) => frag switch
    {
        null => null,
        TSqlScript s => new { kind = "Script", stmts = s.Batches.SelectMany(b => b.Statements).Select(Of).ToList() },
        SelectStatement sel => new { kind = "Select", query = Of(sel.QueryExpression) },
        QuerySpecification qs => new
        {
            kind = "QuerySpec",
            select_list = qs.SelectElements.Select(OfSelectElement).ToList(),
            from = qs.FromClause?.TableReferences.Select(OfTableRef).ToList(),
        },
        BinaryQueryExpression b => new { kind = "BinaryQuery", op = b.BinaryQueryExpressionType.ToString(), left = Of(b.FirstQueryExpression), right = Of(b.SecondQueryExpression) },
        QueryParenthesisExpression p => new { kind = "QueryParen", inner = Of(p.QueryExpression) },
        InsertStatement ins => new { kind = "Insert" },
        _ => new { kind = frag.GetType().Name },
    };

    private static object OfSelectElement(SelectElement e) => e switch
    {
        SelectScalarExpression sse => new
        {
            kind = "SelectScalar",
            alias = sse.ColumnName?.Value,
            expr_kind = sse.Expression?.GetType().Name,
        },
        SelectStarExpression star => new
        {
            kind = "SelectStar",
            qualifier = star.Qualifier != null ? string.Join(".", star.Qualifier.Identifiers.Select(i => i.Value)) : null,
        },
        SelectSetVariable ssv => new
        {
            kind = "SelectSetVariable",
            variable = ssv.Variable?.Name,
            expr_kind = ssv.Expression?.GetType().Name,
        },
        _ => new { kind = e.GetType().Name },
    };

    private static List<string>? NonEmpty(IEnumerable<string>? xs)
    {
        if (xs == null) return null;
        var l = xs.ToList();
        return l.Count == 0 ? null : l;
    }

    private static object OfTableRef(TableReference t)
    {
        var kind = t.GetType().Name;
        string? alias = null;
        List<string>? cols = null;
        List<object>? with_cols = null;

        if (t is TableReferenceWithAlias twa) alias = twa.Alias?.Value;

        switch (t)
        {
            case QueryDerivedTable qdt:
                cols = NonEmpty(qdt.Columns?.Select(c => c.Value));
                break;
            case InlineDerivedTable idt:
                cols = NonEmpty(idt.Columns?.Select(c => c.Value));
                break;
            case SchemaObjectFunctionTableReference sof:
                cols = NonEmpty(sof.Columns?.Select(c => c.Value));
                break;
            case VariableMethodCallTableReference vmc:
                cols = NonEmpty(vmc.Columns?.Select(c => c.Value));
                break;
            case OpenJsonTableReference oj:
                with_cols = oj.SchemaDeclarationItems?.Select(si => (object)new
                {
                    name = si.ColumnDefinition?.ColumnIdentifier?.Value,
                    type = si.ColumnDefinition?.DataType?.Name?.BaseIdentifier?.Value,
                }).ToList();
                break;
            case QualifiedJoin qj:
                return new
                {
                    kind = "QualifiedJoin",
                    join_type = qj.QualifiedJoinType.ToString(),
                    left = OfTableRef(qj.FirstTableReference),
                    right = OfTableRef(qj.SecondTableReference),
                };
            case UnqualifiedJoin uj:
                return new
                {
                    kind = "UnqualifiedJoin",
                    join_type = uj.UnqualifiedJoinType.ToString(),
                    left = OfTableRef(uj.FirstTableReference),
                    right = OfTableRef(uj.SecondTableReference),
                };
        }

        return new
        {
            kind,
            alias,
            columns = cols,
            with_cols,
        };
    }
}
