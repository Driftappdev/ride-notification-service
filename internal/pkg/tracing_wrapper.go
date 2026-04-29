package pkg

import (
	"context"
	gotracing "github.com/driftappdev/libpackage/gotracing"
)

type Tracer = gotracing.Tracer
type TraceSpan = gotracing.Span
type TraceAttrType = gotracing.Attribute

func GlobalTracer(name string) *Tracer { return gotracing.GlobalTracer(name) }
func StartTraceSpan(tracer *Tracer, ctx context.Context, name string) (context.Context, *TraceSpan) {
	return tracer.Start(ctx, name)
}
func TraceAttr(key string, value interface{}) TraceAttrType { return gotracing.Attr(key, value) }
func TraceSetAttrs(span *TraceSpan, attrs ...TraceAttrType) {
	if span != nil {
		span.SetAttribute(attrs...)
	}
}
func TraceSetError(span *TraceSpan, msg string) {
	if span != nil {
		span.SetStatus(gotracing.StatusError, msg)
	}
}
func TraceEnd(span *TraceSpan) {
	if span != nil {
		span.End()
	}
}
