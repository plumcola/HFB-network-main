package uuid

import (
  "github.com/sirupsen/logrus"
  "github.com/google/uuid"
)

type IdHook struct {
  UUID string
}

func (h *IdHook) Levels() []logrus.Level {
  return logrus.AllLevels
}

func (h *IdHook) Fire(entry *logrus.Entry) error {
  entry.Data["uuid"] = h.UUID
  return nil
}

func main() {
  id := uuid.New()
  h := &IdHook{UUID: id.String()}
  logrus.AddHook(h)
  //logrus.Info("info")
  logrus.WithField("uuid", id.String()).Info("test message")
}