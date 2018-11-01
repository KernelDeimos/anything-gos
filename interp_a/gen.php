<?php // code generator because I'm a horrible person

$genMode = "default";
if (count($argv) > 1) $genMode = $argv[1];

$genmsg = array();

if ($genMode != "collapse") {
    $genmsg = array(
        "// -- generated code until ::end (mode=$genMode)"
    );
}

function gentask($file, $name, $f) {
    global $genMode;

    $input = file_get_contents($file);
    $lines = explode("\n", $input);

    $output = array();

    $mode = "default"; // or "generate"

    $genArgs = array();

    foreach ($lines as $i => $line) {
        if ($mode == "default") {
            // Always copy the line to output
            $output[]= $line;

            // Trim line for processing
            $oldLine = $line;
            $line = trim($line);

            // Run generate function if applicable
            $isGenerate = preg_match("/^\/\/::gen/", $line);
            $isShortcut = preg_match("/^\/\/::genl/", $line);
            if ($isGenerate) {
                $parts = explode(' ', $line);

                // Ensure generate directive is what we're looking for
                if ($parts[1] != $name) continue;

                // Get indentation level
                $indent = "";
                $oldLineChars = str_split($oldLine);
                foreach ($oldLineChars as $c) {
                    if ($c == "\t") $indent .= "\t";
                    else break;
                }

                // Generate the code
                $args = array_slice($parts, 2);
                $outputLines = $f($args);
                if ($genMode == "collapse") {
                    $outputLines[0] .= " // -- generated (collapse mode)";
                }
                foreach ($outputLines as $oLine) {
                    $output []= $indent . $oLine;
                }

                if (!$isShortcut) $mode = "ignore";
                else $mode = "ignore-once";
            }
        } else if ($mode == "ignore") {
            $isEnd = preg_match("/^\/\/::end$/", trim($line));
            if ($isEnd) {
                $output[]= $line;
                $mode = "default";
            }
        } else if ($mode == "ignore-once") {
            $mode = "default";
        }
    }

    file_put_contents($file, implode("\n", $output));
}

/*
gentask("idea_a.go", "register-all-functions", function(){
    //
});
*/

$f_verify_args = function($args) {
    global $genmsg, $genMode;

    $lines = $genmsg;

    $fname = $args[0];

    $toCheck = array();

    for ($i = 1; $i < count($args); $i+=2) {
        $varName = $args[$i];
        $varType = $args[$i+1];
        $toCheck[] = array($varName, $varType);
    }

    if ($genMode != "collapse") {
        $lines[]=
            sprintf("if len(args) < %d {", count($toCheck));
        $lines[]=
            sprintf(
                "\treturn nil, errors.New(".
                "\"$fname requires at least %d arguments\")",
                count($toCheck));
        $lines[]=
            sprintf("}\n");
    }

    for ($i=0; $i<count($toCheck); $i++) {
        $item = $toCheck[$i];
        $name = $item[0];
        $type = $item[1];

        $lines[]="var $name $type";
    }

    if ($genMode != "collapse") {
        $lines[]="{";
        $lines[]="\tvar ok bool";
        for ($i=0; $i<count($toCheck); $i++) {
            $item = $toCheck[$i];
            $name = $item[0];
            $type = $item[1];

            $lines[]="\t$name, ok = args[$i].($type)";
            $lines[]="\tif !ok {";
            $lines[]=sprintf("\t\treturn nil, errors.New(".
                "\"%s: argument %d: %s; must be type %s\")",
                $fname, $i, $name, $type);
            $lines[]= "\t}";
        }
        $lines[]="}";
    }

    

    return $lines;
};


$f_errsc_nil_err = function($args) {
    global $genMode;

    if ($genMode == "collapse") return array();
    return array(
        "if err != nil {",
        "\treturn nil, err",
        "}"
    );
};

gentask("builtins_a.go", "verify-args", $f_verify_args);
gentask("evaluator.go", "verify-args", $f_verify_args);

gentask("evaluator.go", "iferr1", $f_errsc_nil_err);
