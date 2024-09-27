Go-Print-Type
=============

The go-print-type project demonstrate how to print a complex type via
reflection.

.. code-block:: bash

    go build text.go
    ./text

`Here <example/output.txt>`_ is an example output where the first field in each
line is the yaml field name.

jsonschema.go demonstrates generating a JSON Schema file for a complex type:

.. code-block:: bash

    go build jsonschema.go
    ./jsonschema
