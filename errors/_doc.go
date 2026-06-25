// Package errors provides structured error creation, wrapping, metadata tags,
// stack capture, and JSON-friendly inspection helpers for Go applications.
//
// Wrap and Wrapf preserve the root Error() text and attach context as metadata.
// Use FormatChain or FullMessage when a joined message is needed for logs.
package errors
